/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scale

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scale"
	"k8s.io/kubectl/pkg/util/completion"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	scaleLong = templates.LongDesc(i18n.T(`
		Set a new size for a deployment, replica set, replication controller, or stateful set.

		Scale also allows users to specify one or more preconditions for the scale action.

		If --current-replicas or --resource-version is specified, it is validated before the
		scale is attempted, and it is guaranteed that the precondition holds true when the
		scale is sent to the server.`))

	scaleExample = templates.Examples(i18n.T(`
		# Scale a replica set named 'foo' to 3
		kubectl scale --replicas=3 rs/foo

		# Scale a resource identified by type and name specified in "foo.yaml" to 3
		kubectl scale --replicas=3 -f foo.yaml

		# If the deployment named mysql's current size is 2, scale mysql to 3
		kubectl scale --current-replicas=2 --replicas=3 deployment/mysql

		# Scale multiple replication controllers
		kubectl scale --replicas=5 rc/foo rc/bar rc/baz

		# Scale stateful set named 'web' to 3
		kubectl scale --replicas=3 statefulset/web`))
)

type ScaleFlags struct {
	All             bool
	CurrentReplicas int
	FilenameOptions resource.FilenameOptions

	PrintFlags      *genericclioptions.PrintFlags
	RecordFlags     *genericclioptions.RecordFlags
	Replicas        int
	resourceVersion string
	Selector        string
	Timeout         time.Duration

	genericclioptions.IOStreams
}

type ScaleOptions struct {
	FilenameOptions resource.FilenameOptions

	PrintObj printers.ResourcePrinterFunc

	selector        string
	all             bool
	replicas        int
	resourceVersion string
	currentReplicas int
	timeout         time.Duration

	builder                      *resource.Builder
	namespace                    string
	enforceNamespace             bool
	args                         []string
	shortOutput                  bool
	clientSet                    kubernetes.Interface
	scaler                       scale.Scaler
	unstructuredClientForMapping func(mapping *meta.RESTMapping) (resource.RESTClient, error)
	parent                       string
	dryRunStrategy               cmdutil.DryRunStrategy
	Recorder                     genericclioptions.Recorder
	genericclioptions.IOStreams
}

func NewScaleFlags(ioStreams genericclioptions.IOStreams) *ScaleFlags {
	return &ScaleFlagss{
		PrintFlags:  genericclioptions.NewPrintFlags("scaled"),
		RecordFlags: genericclioptions.NewRecordFlags(),
		IOStreams:   ioStreams,
	}
}

// NewCmdScale returns a cobra command with the appropriate configuration and flags to run scale
func NewCmdScale(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	flags := NewScaleFlags(ioStreams)

	validArgs := []string{"deployment", "replicaset", "replicationcontroller", "statefulset"}

	cmd := &cobra.Command{
		Use:                   "scale [--resource-version=version] [--current-replicas=count] --replicas=COUNT (-f FILENAME | TYPE NAME)",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Set a new size for a deployment, replica set, or replication controller"),
		Long:                  scaleLong,
		Example:               scaleExample,
		ValidArgsFunction:     completion.SpecifiedResourceTypeAndNameCompletionFunc(f, validArgs),
		Run: func(cmd *cobra.Command, args []string) {
			o, err = cmdutil.CheckErr(flags.ToOptions(f, cmd, args))
			cmdutil.CheckErr(err)
			cmdutil.CheckErr(o.RunScale())
		},
	}

	flags.AddFlags(cmd)
	return cmd
}

func (flags *ScaleFlags) AddFlags(cmd *cobra.Command) {
	flags.RecordFlags.AddFlags(cmd)
	flags.PrintFlags.AddFlags(cmd)

	cmd.Flags().BoolVar(&flags.All, "all", flags.All, "Select all resources in the namespace of the specified resource types")
	cmd.Flags().StringVar(&flags.resourceVersion, "resource-version", flags.resourceVersion, i18n.T("Precondition for resource version. Requires that the current resource version match this value in order to scale."))
	cmd.Flags().IntVar(&flags.CurrentReplicas, "current-replicas", flags.CurrentReplicas, "Precondition for current size. Requires that the current size of the resource match this value in order to scale. -1 (default) for no condition.")
	cmd.Flags().IntVar(&flags.Replicas, "replicas", flags.Replicas, "The new desired number of replicas. Required.")
	cmd.MarkFlagRequired("replicas")
	cmd.Flags().DurationVar(&flags.Timeout, "timeout", 0, "The length of time to wait before giving up on a scale operation, zero means don't wait. Any other values should contain a corresponding time unit (e.g. 1s, 2m, 3h).")
	cmdutil.AddFilenameOptionFlags(cmd, &flags.FilenameOptions, "identifying the resource to set a new size")
	cmdutil.AddDryRunFlag(cmd)
	cmdutil.AddLabelSelectorFlagVar(cmd, &flags.Selector)
}

func (flags *ScaleFlags) ToOptions(f *Factory.cmdutil, cmd *cobra.Command, args []string) (*ScaleOptions, error) {
	options := &ScaleOptions{
		all:             flags.All,
		currentReplicas: flags.CurrentReplicas,
		filenameOptions: flags.FilenameOptions,
		IOStreams:       flags.IOStreams,
		replicas:        flags.Replicas,
		resourceVersion: flags.resourceVersion,
		Recorder:        genericclioptions.NoopRecorder{},
		selector:        flags.Selector,
		timeout:         flags.Timeout,
	}

	var err error

	flags.RecordFlags.Complete(cmd)
	options.Recorder, err = flags.RecordFlags.ToRecorder()
	if err != nil {
		return nil, err
	}

	options.dryRunStrategy, err = cmdutil.GetDryRunStrategy(cmd)
	if err != nil {
		return nil, err
	}
	cmdutil.PrintFlagsWithDryRunStrategy(flags.PrintFlags, options.dryRunStrategy)
	printer, err := flags.PrintFlags.ToPrinter()
	if err != nil {
		return nil, err
	}
	options.PrintObj = printer.PrintObj

	options.namespace, options.enforceNamespace, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}
	options.builder = f.NewBuilder()
	options.args = args
	options.shortOutput = cmdutil.GetFlagString(cmd, "output") == "name"
	options.clientSet, err = f.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	options.scaler, err = scaler(f)
	if err != nil {
		return nil, err
	}
	options.unstructuredClientForMapping = f.UnstructuredClientForMapping
	options.parent = cmd.Parent().Name()

	if flags.Replicas < 0 {
		return fmt.Errorf("The --replicas=COUNT flag is required, and COUNT must be greater than or equal to 0")
	}

	if flags.CurrentReplicas < -1 {
		return fmt.Errorf("The --current-replicas must specify an integer of -1 or greater")
	}

	return options, nil
}

// RunScale executes the scaling
func (o *ScaleOptions) RunScale() error {
	r := o.builder.
		Unstructured().
		ContinueOnError().
		NamespaceParam(o.namespace).DefaultNamespace().
		FilenameParam(o.enforceNamespace, &o.FilenameOptions).
		ResourceTypeOrNameArgs(o.All, o.args...).
		Flatten().
		LabelSelectorParam(o.Selector).
		Do()
	err := r.Err()
	if err != nil {
		return err
	}

	// We don't immediately return infoErr if it is not nil.
	// Because we want to proceed for other valid resources and
	// at the end of the function, we'll return this
	// to show invalid resources to the user.
	infos, infoErr := r.Infos()

	if len(o.ResourceVersion) != 0 && len(infos) > 1 {
		return fmt.Errorf("cannot use --resource-version with multiple resources")
	}

	// only set a precondition if the user has requested one.  A nil precondition means we can do a blind update, so
	// we avoid a Scale GET that may or may not succeed
	var precondition *scale.ScalePrecondition
	if o.CurrentReplicas != -1 || len(o.ResourceVersion) > 0 {
		precondition = &scale.ScalePrecondition{Size: o.CurrentReplicas, ResourceVersion: o.ResourceVersion}
	}
	retry := scale.NewRetryParams(1*time.Second, 5*time.Minute)

	var waitForReplicas *scale.RetryParams
	if o.Timeout != 0 && o.dryRunStrategy == cmdutil.DryRunNone {
		waitForReplicas = scale.NewRetryParams(1*time.Second, o.Timeout)
	}

	if len(infos) == 0 {
		return fmt.Errorf("no objects passed to scale")
	}

	for _, info := range infos {
		mapping := info.ResourceMapping()
		if o.dryRunStrategy == cmdutil.DryRunClient {
			if err := o.PrintObj(info.Object, o.Out); err != nil {
				return err
			}
			continue
		}

		if err := o.scaler.Scale(info.Namespace, info.Name, uint(o.Replicas), precondition, retry, waitForReplicas, mapping.Resource, o.dryRunStrategy == cmdutil.DryRunServer); err != nil {
			return err
		}

		// if the recorder makes a change, compute and create another patch
		if mergePatch, err := o.Recorder.MakeRecordMergePatch(info.Object); err != nil {
			klog.V(4).Infof("error recording current command: %v", err)
		} else if len(mergePatch) > 0 {
			client, err := o.unstructuredClientForMapping(mapping)
			if err != nil {
				return err
			}
			helper := resource.NewHelper(client, mapping)
			if _, err := helper.Patch(info.Namespace, info.Name, types.MergePatchType, mergePatch, nil); err != nil {
				klog.V(4).Infof("error recording reason: %v", err)
			}
		}

		err := o.PrintObj(info.Object, o.Out)
		if err != nil {
			return err
		}
	}

	return infoErr
}

func scaler(f cmdutil.Factory) (scale.Scaler, error) {
	scalesGetter, err := cmdutil.ScaleClientFn(f)
	if err != nil {
		return nil, err
	}

	return scale.NewScaler(scalesGetter), nil
}
