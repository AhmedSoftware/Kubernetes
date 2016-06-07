/*
Copyright 2016 The Kubernetes Authors.

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

package flexvolume

import (
	"path"
	"strings"
	"sync"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/types"
	"k8s.io/kubernetes/pkg/util/exec"
	"k8s.io/kubernetes/pkg/util/mount"
	utilstrings "k8s.io/kubernetes/pkg/util/strings"
	"k8s.io/kubernetes/pkg/volume"
)

const flexVolumePluginName = "kubernetes.io/flexvolume"

// FlexVolumePlugin object.
type flexVolumePlugin struct {
	driverName string
	execPath   string
	host       volume.VolumeHost
	runner     exec.Interface

	sync.Mutex
	unsupportedCommands []string
}

var _ volume.AttachableVolumePlugin = &flexVolumePlugin{}
var _ volume.PersistentVolumePlugin = &flexVolumePlugin{}

// Init is part of the volume.VolumePlugin interface.
func (plugin *flexVolumePlugin) Init(host volume.VolumeHost) error {
	plugin.host = host
	// call the init script
	call := plugin.NewDriverCall(initCmd)
	_, err := call.Run()
	return err
}

func (plugin *flexVolumePlugin) getExecutable() string {
	parts := strings.Split(plugin.driverName, "/")
	execName := parts[len(parts)-1]
	return path.Join(plugin.execPath, execName)
}

// Name is part of the volume.VolumePlugin interface.
func (plugin *flexVolumePlugin) GetPluginName() string {
	return plugin.driverName
}

// GetVolumeName is part of the volume.VolumePlugin interface.
func (plugin *flexVolumePlugin) GetVolumeName(spec *volume.Spec) (string, error) {
	call := plugin.NewDriverCall(getVolumeNameCmd)
	call.AppendSpec(spec, plugin.host, nil)

	status, err := call.Run()
	if isCmdNotSupportedErr(err) {
		return (*pluginDefaults)(plugin).GetVolumeName(spec)
	} else if err != nil {
		return "", err
	}
	return utilstrings.EscapeQualifiedNameForDisk(status.Device), nil
}

// CanSupport is part of the volume.VolumePlugin interface.
func (plugin *flexVolumePlugin) CanSupport(spec *volume.Spec) bool {
	source, _ := getVolumeSource(spec)
	return (source != nil) && (source.Driver == plugin.driverName)
}

// RequiresRemount is part of the volume.VolumePlugin interface.
func (plugin *flexVolumePlugin) RequiresRemount() bool {
	return false
}

// GetAccessModes gets the allowed access modes for this plugin.
func (plugin *flexVolumePlugin) GetAccessModes() []api.PersistentVolumeAccessMode {
	return []api.PersistentVolumeAccessMode{
		api.ReadWriteOnce,
		api.ReadOnlyMany,
	}
}

// NewMounter is part of the volume.VolumePlugin interface.
func (plugin *flexVolumePlugin) NewMounter(spec *volume.Spec, pod *api.Pod, _ volume.VolumeOptions) (volume.Mounter, error) {
	secrets, err := getSecrets(spec, plugin.host)
	if err != nil {
		return nil, err
	}
	return plugin.newMounterInternal(spec, pod, plugin.host.GetMounter(), plugin.runner, secrets)
}

// newMounterInternal is the internal mounter routine to build the volume.
func (plugin *flexVolumePlugin) newMounterInternal(spec *volume.Spec, pod *api.Pod, mounter mount.Interface, runner exec.Interface, secrets map[string]string) (volume.Mounter, error) {
	source, readOnly := getVolumeSource(spec)
	return &flexVolumeMounter{
		flexVolume: &flexVolume{
			driverName: source.Driver,
			execPath:   plugin.getExecutable(),
			mounter:    mounter,
			plugin:     plugin,
			podUID:     pod.UID,
			volName:    spec.Name(),
		},
		runner:             runner,
		spec:               spec,
		readOnly:           readOnly,
		blockDeviceMounter: &mount.SafeFormatAndMount{Interface: mounter, Runner: runner},
	}, nil
}

// NewUnmounter is part of the volume.VolumePlugin interface.
func (plugin *flexVolumePlugin) NewUnmounter(volName string, podUID types.UID) (volume.Unmounter, error) {
	return plugin.newUnmounterInternal(volName, podUID, plugin.host.GetMounter(), plugin.runner)
}

// newUnmounterInternal is the internal unmounter routine to clean the volume.
func (plugin *flexVolumePlugin) newUnmounterInternal(volName string, podUID types.UID, mounter mount.Interface, runner exec.Interface) (volume.Unmounter, error) {
	return &flexVolumeUnmounter{
		flexVolume: &flexVolume{
			driverName: plugin.driverName,
			execPath:   plugin.getExecutable(),
			mounter:    mounter,
			plugin:     plugin,
			podUID:     podUID,
			volName:    volName,
		},
		runner: runner,
	}, nil
}

// NewAttacher is part of the volume.AttachableVolumePlugin interface.
func (plugin *flexVolumePlugin) NewAttacher() (volume.Attacher, error) {
	return &flexVolumeAttacher{plugin}, nil
}

// NewDetacher is part of the volume.AttachableVolumePlugin interface.
func (plugin *flexVolumePlugin) NewDetacher() (volume.Detacher, error) {
	return &flexVolumeDetacher{plugin}, nil
}

// ConstructVolumeSpec is part of the volume.AttachableVolumePlugin interface.
func (plugin *flexVolumePlugin) ConstructVolumeSpec(volumeName, mountPath string) (*volume.Spec, error) {
	flexVolume := &api.Volume{
		Name: volumeName,
		VolumeSource: api.VolumeSource{
			FlexVolume: &api.FlexVolumeSource{
				Driver: plugin.driverName,
			},
		},
	}
	// namespace is unknown here
	return volume.NewSpecFromVolume(flexVolume, ""), nil
}

// Mark the given commands as unsupported.
func (plugin *flexVolumePlugin) unsupported(commands ...string) {
	plugin.Lock()
	defer plugin.Unlock()
	plugin.unsupportedCommands = append(plugin.unsupportedCommands, commands...)
}

// Returns true iff the given command is know to be unsupported.
func (plugin *flexVolumePlugin) isUnsupported(command string) bool {
	plugin.Lock()
	defer plugin.Unlock()
	for _, unsupportedCommand := range plugin.unsupportedCommands {
		if command == unsupportedCommand {
			return true
		}
	}
	return false
}

func (plugin *flexVolumePlugin) GetDeviceMountRefs(deviceMountPath string) ([]string, error) {
	mounter := plugin.host.GetMounter()
	return mount.GetMountRefs(mounter, deviceMountPath)
}
