/*
Copyright 2024.

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BackupPhase string

const (
	BackupPhasePending         BackupPhase = "Pending"
	BackupPhaseCapturing       BackupPhase = "Capturing"
	BackupPhaseWritingManifest BackupPhase = "WritingManifest"
	BackupPhaseSucceeded       BackupPhase = "Succeeded"
	BackupPhaseFailed          BackupPhase = "Failed"
)

// PlatformBackupSpec defines the desired state of PlatformBackup
type PlatformBackupSpec struct {
	Storage    StorageSpec    `json:"storage"`
	Components ComponentsSpec `json:"components"`
}

type StorageSpec struct {
	S3 S3StorageSpec `json:"s3"`
}

type S3StorageSpec struct {
	Endpoint       string                    `json:"endpoint"`
	Bucket         string                    `json:"bucket"`
	Region         string                    `json:"region,omitempty"`
	CredentialsRef corev1.LocalObjectReference `json:"credentialsRef"`
}

type ComponentsSpec struct {
	Etcd   ComponentToggle `json:"etcd,omitempty"`
	CNPG   ComponentToggle `json:"cnpg,omitempty"`
	Velero ComponentToggle `json:"velero,omitempty"`
}

type ComponentToggle struct {
	Enabled bool `json:"enabled"`
}

// PlatformBackupStatus defines the observed state of PlatformBackup
type PlatformBackupStatus struct {
	Phase              BackupPhase      `json:"phase,omitempty"`
	BackupID           string           `json:"backupID,omitempty"`
	TopologyDigest     string           `json:"topologyDigest,omitempty"`
	Artefacts          ArtefactsStatus  `json:"artefacts,omitempty"`
	ObservedGeneration int64            `json:"observedGeneration,omitempty"`
	NextReconcileTime  metav1.Time      `json:"nextReconcileTime,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

type ArtefactsStatus struct {
	Etcd   *EtcdArtefact   `json:"etcd,omitempty"`
	CNPG   *CNPGArtefact   `json:"cnpg,omitempty"`
	Velero *VeleroArtefact `json:"velero,omitempty"`
}

type EtcdArtefact struct {
	Shards map[string]EtcdShardArtefact `json:"shards,omitempty"`
}

type EtcdShardArtefact struct {
	SnapshotKey  string      `json:"snapshotKey"`
	SnapshotTime metav1.Time `json:"snapshotTime"`
}

type CNPGArtefact struct {
	Backups map[string]string `json:"backups,omitempty"`
}

type VeleroArtefact struct {
	BackupName string `json:"backupName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=.status.phase
// +kubebuilder:printcolumn:name="BackupID",type=string,JSONPath=.status.backupID
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`

// PlatformBackup is the Schema for the platformbackups API
type PlatformBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlatformBackupSpec   `json:"spec,omitempty"`
	Status PlatformBackupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PlatformBackupList contains a list of PlatformBackup
type PlatformBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PlatformBackup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PlatformBackup{}, &PlatformBackupList{})
}

func (b *PlatformBackup) GetConditions() []metav1.Condition           { return b.Status.Conditions }
func (b *PlatformBackup) SetConditions(c []metav1.Condition)          { b.Status.Conditions = c }
func (b *PlatformBackup) GetObservedGeneration() int64                { return b.Status.ObservedGeneration }
func (b *PlatformBackup) SetObservedGeneration(g int64)               { b.Status.ObservedGeneration = g }
func (b *PlatformBackup) GetNextReconcileTime() metav1.Time           { return b.Status.NextReconcileTime }
func (b *PlatformBackup) SetNextReconcileTime(t metav1.Time)          { b.Status.NextReconcileTime = t }
