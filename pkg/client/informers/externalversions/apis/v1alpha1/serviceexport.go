/*
Copyright 2020 The Kubernetes Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	pkgapisv1alpha1 "sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"
	versioned "sigs.k8s.io/mcs-api/pkg/client/clientset/versioned"
	internalinterfaces "sigs.k8s.io/mcs-api/pkg/client/informers/externalversions/internalinterfaces"
	apisv1alpha1 "sigs.k8s.io/mcs-api/pkg/client/listers/apis/v1alpha1"
)

// ServiceExportInformer provides access to a shared informer and lister for
// ServiceExports.
type ServiceExportInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() apisv1alpha1.ServiceExportLister
}

type serviceExportInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewServiceExportInformer constructs a new informer for ServiceExport type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewServiceExportInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredServiceExportInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredServiceExportInformer constructs a new informer for ServiceExport type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredServiceExportInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MulticlusterV1alpha1().ServiceExports(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MulticlusterV1alpha1().ServiceExports(namespace).Watch(context.TODO(), options)
			},
		},
		&pkgapisv1alpha1.ServiceExport{},
		resyncPeriod,
		indexers,
	)
}

func (f *serviceExportInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredServiceExportInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *serviceExportInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&pkgapisv1alpha1.ServiceExport{}, f.defaultInformer)
}

func (f *serviceExportInformer) Lister() apisv1alpha1.ServiceExportLister {
	return apisv1alpha1.NewServiceExportLister(f.Informer().GetIndexer())
}
