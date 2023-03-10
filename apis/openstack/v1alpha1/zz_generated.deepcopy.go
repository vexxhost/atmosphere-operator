//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"encoding/json"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Barbican) DeepCopyInto(out *Barbican) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Barbican.
func (in *Barbican) DeepCopy() *Barbican {
	if in == nil {
		return nil
	}
	out := new(Barbican)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Barbican) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BarbicanList) DeepCopyInto(out *BarbicanList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Barbican, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BarbicanList.
func (in *BarbicanList) DeepCopy() *BarbicanList {
	if in == nil {
		return nil
	}
	out := new(BarbicanList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BarbicanList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BarbicanSpec) DeepCopyInto(out *BarbicanSpec) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.KeystoneRef = in.KeystoneRef
	out.SecretsRef = in.SecretsRef
	out.RabbitmqReference = in.RabbitmqReference
	out.DatabaseReference = in.DatabaseReference
	in.Overrides.DeepCopyInto(&out.Overrides)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BarbicanSpec.
func (in *BarbicanSpec) DeepCopy() *BarbicanSpec {
	if in == nil {
		return nil
	}
	out := new(BarbicanSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BarbicanStatus) DeepCopyInto(out *BarbicanStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BarbicanStatus.
func (in *BarbicanStatus) DeepCopy() *BarbicanStatus {
	if in == nil {
		return nil
	}
	out := new(BarbicanStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Designate) DeepCopyInto(out *Designate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Designate.
func (in *Designate) DeepCopy() *Designate {
	if in == nil {
		return nil
	}
	out := new(Designate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Designate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DesignateList) DeepCopyInto(out *DesignateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Designate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DesignateList.
func (in *DesignateList) DeepCopy() *DesignateList {
	if in == nil {
		return nil
	}
	out := new(DesignateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DesignateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DesignateSpec) DeepCopyInto(out *DesignateSpec) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.KeystoneRef = in.KeystoneRef
	out.SecretsRef = in.SecretsRef
	out.RabbitmqReference = in.RabbitmqReference
	out.DatabaseReference = in.DatabaseReference
	in.Overrides.DeepCopyInto(&out.Overrides)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DesignateSpec.
func (in *DesignateSpec) DeepCopy() *DesignateSpec {
	if in == nil {
		return nil
	}
	out := new(DesignateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DesignateStatus) DeepCopyInto(out *DesignateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DesignateStatus.
func (in *DesignateStatus) DeepCopy() *DesignateStatus {
	if in == nil {
		return nil
	}
	out := new(DesignateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Glance) DeepCopyInto(out *Glance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Glance.
func (in *Glance) DeepCopy() *Glance {
	if in == nil {
		return nil
	}
	out := new(Glance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Glance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlanceList) DeepCopyInto(out *GlanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Glance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlanceList.
func (in *GlanceList) DeepCopy() *GlanceList {
	if in == nil {
		return nil
	}
	out := new(GlanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GlanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlanceSpec) DeepCopyInto(out *GlanceSpec) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.KeystoneRef = in.KeystoneRef
	out.HorizonRef = in.HorizonRef
	out.SecretsRef = in.SecretsRef
	out.RabbitmqReference = in.RabbitmqReference
	out.DatabaseReference = in.DatabaseReference
	in.Overrides.DeepCopyInto(&out.Overrides)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlanceSpec.
func (in *GlanceSpec) DeepCopy() *GlanceSpec {
	if in == nil {
		return nil
	}
	out := new(GlanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlanceStatus) DeepCopyInto(out *GlanceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlanceStatus.
func (in *GlanceStatus) DeepCopy() *GlanceStatus {
	if in == nil {
		return nil
	}
	out := new(GlanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelmOverrides) DeepCopyInto(out *HelmOverrides) {
	*out = *in
	if in.RawMessage != nil {
		in, out := &in.RawMessage, &out.RawMessage
		*out = make(json.RawMessage, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelmOverrides.
func (in *HelmOverrides) DeepCopy() *HelmOverrides {
	if in == nil {
		return nil
	}
	out := new(HelmOverrides)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Horizon) DeepCopyInto(out *Horizon) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Horizon.
func (in *Horizon) DeepCopy() *Horizon {
	if in == nil {
		return nil
	}
	out := new(Horizon)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Horizon) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizonList) DeepCopyInto(out *HorizonList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Horizon, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizonList.
func (in *HorizonList) DeepCopy() *HorizonList {
	if in == nil {
		return nil
	}
	out := new(HorizonList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HorizonList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizonSpec) DeepCopyInto(out *HorizonSpec) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.KeystoneRef = in.KeystoneRef
	out.DatabaseReference = in.DatabaseReference
	out.SecretsRef = in.SecretsRef
	in.Overrides.DeepCopyInto(&out.Overrides)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizonSpec.
func (in *HorizonSpec) DeepCopy() *HorizonSpec {
	if in == nil {
		return nil
	}
	out := new(HorizonSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizonStatus) DeepCopyInto(out *HorizonStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizonStatus.
func (in *HorizonStatus) DeepCopy() *HorizonStatus {
	if in == nil {
		return nil
	}
	out := new(HorizonStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressConfig) DeepCopyInto(out *IngressConfig) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.TLS.DeepCopyInto(&out.TLS)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressConfig.
func (in *IngressConfig) DeepCopy() *IngressConfig {
	if in == nil {
		return nil
	}
	out := new(IngressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ironic) DeepCopyInto(out *Ironic) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ironic.
func (in *Ironic) DeepCopy() *Ironic {
	if in == nil {
		return nil
	}
	out := new(Ironic)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Ironic) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IronicList) DeepCopyInto(out *IronicList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Ironic, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IronicList.
func (in *IronicList) DeepCopy() *IronicList {
	if in == nil {
		return nil
	}
	out := new(IronicList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IronicList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IronicSpec) DeepCopyInto(out *IronicSpec) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.KeystoneRef = in.KeystoneRef
	out.SecretsRef = in.SecretsRef
	out.RabbitmqReference = in.RabbitmqReference
	out.DatabaseReference = in.DatabaseReference
	in.Overrides.DeepCopyInto(&out.Overrides)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IronicSpec.
func (in *IronicSpec) DeepCopy() *IronicSpec {
	if in == nil {
		return nil
	}
	out := new(IronicSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IronicStatus) DeepCopyInto(out *IronicStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IronicStatus.
func (in *IronicStatus) DeepCopy() *IronicStatus {
	if in == nil {
		return nil
	}
	out := new(IronicStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Keystone) DeepCopyInto(out *Keystone) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Keystone.
func (in *Keystone) DeepCopy() *Keystone {
	if in == nil {
		return nil
	}
	out := new(Keystone)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Keystone) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeystoneList) DeepCopyInto(out *KeystoneList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Keystone, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeystoneList.
func (in *KeystoneList) DeepCopy() *KeystoneList {
	if in == nil {
		return nil
	}
	out := new(KeystoneList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KeystoneList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeystoneSpec) DeepCopyInto(out *KeystoneSpec) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.SecretsRef = in.SecretsRef
	out.RabbitmqReference = in.RabbitmqReference
	out.DatabaseReference = in.DatabaseReference
	in.Overrides.DeepCopyInto(&out.Overrides)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeystoneSpec.
func (in *KeystoneSpec) DeepCopy() *KeystoneSpec {
	if in == nil {
		return nil
	}
	out := new(KeystoneSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeystoneStatus) DeepCopyInto(out *KeystoneStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeystoneStatus.
func (in *KeystoneStatus) DeepCopy() *KeystoneStatus {
	if in == nil {
		return nil
	}
	out := new(KeystoneStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacedName) DeepCopyInto(out *NamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacedName.
func (in *NamespacedName) DeepCopy() *NamespacedName {
	if in == nil {
		return nil
	}
	out := new(NamespacedName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Neutron) DeepCopyInto(out *Neutron) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Neutron.
func (in *Neutron) DeepCopy() *Neutron {
	if in == nil {
		return nil
	}
	out := new(Neutron)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Neutron) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NeutronList) DeepCopyInto(out *NeutronList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Neutron, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NeutronList.
func (in *NeutronList) DeepCopy() *NeutronList {
	if in == nil {
		return nil
	}
	out := new(NeutronList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NeutronList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NeutronSpec) DeepCopyInto(out *NeutronSpec) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.KeystoneRef = in.KeystoneRef
	out.NovaRef = in.NovaRef
	out.OctaviaRef = in.OctaviaRef
	out.DesignateRef = in.DesignateRef
	out.IronicRef = in.IronicRef
	out.CoreDNSRef = in.CoreDNSRef
	out.SecretsRef = in.SecretsRef
	out.RabbitmqReference = in.RabbitmqReference
	out.DatabaseReference = in.DatabaseReference
	in.Overrides.DeepCopyInto(&out.Overrides)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NeutronSpec.
func (in *NeutronSpec) DeepCopy() *NeutronSpec {
	if in == nil {
		return nil
	}
	out := new(NeutronSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NeutronStatus) DeepCopyInto(out *NeutronStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NeutronStatus.
func (in *NeutronStatus) DeepCopy() *NeutronStatus {
	if in == nil {
		return nil
	}
	out := new(NeutronStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Nova) DeepCopyInto(out *Nova) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Nova.
func (in *Nova) DeepCopy() *Nova {
	if in == nil {
		return nil
	}
	out := new(Nova)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Nova) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaList) DeepCopyInto(out *NovaList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Nova, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaList.
func (in *NovaList) DeepCopy() *NovaList {
	if in == nil {
		return nil
	}
	out := new(NovaList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaSpec) DeepCopyInto(out *NovaSpec) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	in.VncIngress.DeepCopyInto(&out.VncIngress)
	out.KeystoneRef = in.KeystoneRef
	out.PlacementRef = in.PlacementRef
	out.GlanceRef = in.GlanceRef
	out.NeutronRef = in.NeutronRef
	out.IronicRef = in.IronicRef
	out.SecretsRef = in.SecretsRef
	out.RabbitmqReference = in.RabbitmqReference
	out.DatabaseReference = in.DatabaseReference
	in.Overrides.DeepCopyInto(&out.Overrides)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaSpec.
func (in *NovaSpec) DeepCopy() *NovaSpec {
	if in == nil {
		return nil
	}
	out := new(NovaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaStatus) DeepCopyInto(out *NovaStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaStatus.
func (in *NovaStatus) DeepCopy() *NovaStatus {
	if in == nil {
		return nil
	}
	out := new(NovaStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Octavia) DeepCopyInto(out *Octavia) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Octavia.
func (in *Octavia) DeepCopy() *Octavia {
	if in == nil {
		return nil
	}
	out := new(Octavia)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Octavia) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OctaviaAmphoraConfig) DeepCopyInto(out *OctaviaAmphoraConfig) {
	*out = *in
	out.ServerCertificateAuthority = in.ServerCertificateAuthority
	out.ClientCertificate = in.ClientCertificate
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OctaviaAmphoraConfig.
func (in *OctaviaAmphoraConfig) DeepCopy() *OctaviaAmphoraConfig {
	if in == nil {
		return nil
	}
	out := new(OctaviaAmphoraConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OctaviaList) DeepCopyInto(out *OctaviaList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Octavia, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OctaviaList.
func (in *OctaviaList) DeepCopy() *OctaviaList {
	if in == nil {
		return nil
	}
	out := new(OctaviaList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OctaviaList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OctaviaSpec) DeepCopyInto(out *OctaviaSpec) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.KeystoneRef = in.KeystoneRef
	out.NeutronRef = in.NeutronRef
	out.SecretsRef = in.SecretsRef
	out.RabbitmqReference = in.RabbitmqReference
	out.DatabaseReference = in.DatabaseReference
	out.AmphoraConfig = in.AmphoraConfig
	if in.HealthManagers != nil {
		in, out := &in.HealthManagers, &out.HealthManagers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Overrides.DeepCopyInto(&out.Overrides)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OctaviaSpec.
func (in *OctaviaSpec) DeepCopy() *OctaviaSpec {
	if in == nil {
		return nil
	}
	out := new(OctaviaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OctaviaStatus) DeepCopyInto(out *OctaviaStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OctaviaStatus.
func (in *OctaviaStatus) DeepCopy() *OctaviaStatus {
	if in == nil {
		return nil
	}
	out := new(OctaviaStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Placement) DeepCopyInto(out *Placement) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Placement.
func (in *Placement) DeepCopy() *Placement {
	if in == nil {
		return nil
	}
	out := new(Placement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Placement) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementList) DeepCopyInto(out *PlacementList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Placement, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementList.
func (in *PlacementList) DeepCopy() *PlacementList {
	if in == nil {
		return nil
	}
	out := new(PlacementList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlacementList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementSpec) DeepCopyInto(out *PlacementSpec) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.KeystoneRef = in.KeystoneRef
	out.SecretsRef = in.SecretsRef
	out.DatabaseReference = in.DatabaseReference
	in.Overrides.DeepCopyInto(&out.Overrides)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementSpec.
func (in *PlacementSpec) DeepCopy() *PlacementSpec {
	if in == nil {
		return nil
	}
	out := new(PlacementSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementStatus) DeepCopyInto(out *PlacementStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementStatus.
func (in *PlacementStatus) DeepCopy() *PlacementStatus {
	if in == nil {
		return nil
	}
	out := new(PlacementStatus)
	in.DeepCopyInto(out)
	return out
}
