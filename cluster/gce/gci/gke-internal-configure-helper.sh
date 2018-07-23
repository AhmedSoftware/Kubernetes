#!/bin/bash

function start_nap {
  if [[ "${ENABLE_NAP:-}" == "true" ]]; then
    echo "Start Node Auto-Provisioning (NAP)"
    local -r manifests_dir="${KUBE_HOME}/kube-manifests/kubernetes/gci-trusty"

    # Re-using Cluster Autoscaler setup functions from OSS
    setup-addon-manifests "addons" "rbac/cluster-autoscaler"
    create-clusterautoscaler-kubeconfig
    prepare-log-file /var/log/cluster-autoscaler.log

    # Add our GKE specific CRD
    mkdir -p "${manifests_dir}/autoscaling"
    cp "${manifests_dir}/internal-capacity-request-crd.yaml" "${manifests_dir}/autoscaling"
    setup-addon-manifests "addons" "autoscaling"

    # Prepare Autoscaler manifest
    local -r src_file="${manifests_dir}/internal-cluster-autoscaler.manifest"
    local params="${CLOUD_CONFIG_OPT} ${NAP_CONFIG:-}"

    sed -i -e "s@{{params}}@${params:-}@g" "${src_file}"
    sed -i -e "s@{{cloud_config_mount}}@${CLOUD_CONFIG_MOUNT}@g" "${src_file}"
    sed -i -e "s@{{cloud_config_volume}}@${CLOUD_CONFIG_VOLUME}@g" "${src_file}"
    sed -i -e "s@{%.*%}@@g" "${src_file}"

    cp "${src_file}" /etc/kubernetes/manifests
  fi
}

function start_vertical_pod_autoscaler {
  if [[ "${ENABLE_VERTICAL_POD_AUTOSCALER:-}" == "true" ]]; then
    echo "Start Vertical Pod Autoscaler (VPA)"
    generate_vertical_pod_autoscaler_admission_controller_certs

    local -r manifests_dir="${KUBE_HOME}/kube-manifests/kubernetes/gci-trusty"

    mkdir -p "${manifests_dir}/vertical-pod-autoscaler"

    cp "${manifests_dir}/internal-vpa-crd.yaml" "${manifests_dir}/vertical-pod-autoscaler"
    cp "${manifests_dir}/internal-vpa-rbac.yaml" "${manifests_dir}/vertical-pod-autoscaler"
    setup-addon-manifests "addons" "vertical-pod-autoscaler"

    for component in admission-controller recommender updater; do
      prepare-log-file /var/log/vpa-${component}.log

      # Prepare manifest
      local src_file="${manifests_dir}/internal-vpa-${component}.manifest"

      sed -i -e "s@{{cloud_config_mount}}@${CLOUD_CONFIG_MOUNT}@g" "${src_file}"
      sed -i -e "s@{{cloud_config_volume}}@${CLOUD_CONFIG_VOLUME}@g" "${src_file}"
      sed -i -e "s@{%.*%}@@g" "${src_file}"

      cp "${src_file}" /etc/kubernetes/manifests
    done

  fi
}

function generate_vertical_pod_autoscaler_admission_controller_certs {
  local certs_dir="/etc/tls-certs" #TODO: what is the best place for self-singed certs?
  echo "Generating certs for the VPA Admission Controller in ${certs_dir}."
  mkdir -p ${certs_dir}
  cat > ${certs_dir}/server.conf << EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
EOF

  # Create a certificate authority
  openssl genrsa -out ${certs_dir}/caKey.pem 2048
  openssl req -x509 -new -nodes -key ${certs_dir}/caKey.pem -days 100000 -out ${certs_dir}/caCert.pem -subj "/CN=vpa_webhook_ca"

  # Create a server certiticate
  openssl genrsa -out ${certs_dir}/serverKey.pem 2048
  # Note the CN is the DNS name of the service of the webhook.
  # TODO(b/111244006) For now admission controller is running as localhost
  openssl req -new -key ${certs_dir}/serverKey.pem -out ${certs_dir}/server.csr -subj "/CN=localhost" -config ${certs_dir}/server.conf
  openssl x509 -req -in ${certs_dir}/server.csr -CA ${certs_dir}/caCert.pem -CAkey ${certs_dir}/caKey.pem -CAcreateserial -out ${certs_dir}/serverCert.pem -days 100000 -extensions v3_req -extfile ${certs_dir}/server.conf
}


function gke-internal-master-start {
  echo "Internal GKE configuration start"
  compute-master-manifest-variables
  start_nap
  start_vertical_pod_autoscaler
  echo "Internal GKE configuration done"
}
