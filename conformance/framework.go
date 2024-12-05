/*
Copyright 2023 The Kubernetes Authors.

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

package conformance

import (
	"bytes"
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type (
	DoOperationFunc func(ctx context.Context) (interface{}, error)
	CheckResultFunc func(result interface{}) (bool, string, error)
)

// AwaitResultOrError periodically performs the given operation until the given CheckResultFunc returns true, an error, or a
// timeout is reached.
func AwaitResultOrError(opMsg string, doOperation DoOperationFunc, checkResult CheckResultFunc) (interface{}, error) {
	var finalResult interface{}
	var lastMsg string

	err := wait.PollUntilContextTimeout(context.Background(), time.Second,
		20*time.Second, true, func(ctx context.Context) (bool, error) {
			result, err := doOperation(ctx)
			if err != nil {
				if IsTransientError(err, opMsg) {
					return false, nil
				}
				return false, err
			}

			ok, msg, err := checkResult(result)
			if err != nil {
				return false, err
			}

			if ok {
				finalResult = result
				return true, nil
			}

			lastMsg = msg
			return false, nil
		})

	if err != nil {
		errMsg := "Failed to " + opMsg
		if lastMsg != "" {
			errMsg += ". " + lastMsg
		}

		if wait.Interrupted(err) {
			err = errors.New(errMsg)
		} else {
			err = errors.Wrap(err, errMsg)
		}
	}

	return finalResult, err
}

// identify API errors which could be considered transient/recoverable
// due to server state.
func IsTransientError(err error, opMsg string) bool {
	if apierrors.IsInternalError(err) ||
		apierrors.IsServerTimeout(err) ||
		apierrors.IsTimeout(err) ||
		apierrors.IsServiceUnavailable(err) ||
		apierrors.IsUnexpectedServerError(err) ||
		apierrors.IsTooManyRequests(err) {
		fmt.Fprintf(GinkgoWriter, "Transient failure when attempting to %s: %v", opMsg, err)
		return true
	}

	return false
}

func execCmd(k8s kubernetes.Interface, config *rest.Config, podName string, podNamespace string, command []string) ([]byte, []byte, error) {
	req := k8s.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(podNamespace).SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Command: command,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return []byte{}, []byte{}, err
	}

	var stdout, stderr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return stdout.Bytes(), stderr.Bytes(), nil
}
