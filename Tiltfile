# -*- mode: Python -*-

load('ext://restart_process', 'docker_build_with_restart')

allow_k8s_contexts([
  'rancher-desktop',
  'minikube',
  'dockerdesktop',
  'kind',
])

compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chaos ./'

local_resource(
  'go-compile',
  compile_cmd,
  deps=[
    './main.go',
    './pkg',
  ],
)

docker_build_with_restart(
  'podchaosmonkey',
  '.',
  entrypoint=['/app/chaos'],
  dockerfile='deployment/Dockerfile',
  # only=[
  #   './chaos',
  #   './production.env'
  # ],
  live_update=[
    sync('./chaos', '/app/chaos'),
    sync('./production.env', '/app/production.env'),
  ],
)

k8s_yaml('deployment/kubernetes.yaml')
k8s_resource(
  'podchaosmonkey',
  objects=[
    'podchaosmonkey:serviceaccount',
    'podchaosmonkey:clusterrole',
    'podchaosmonkey:clusterrolebinding',
    'choasmonkey:namespace',
  ],
  resource_deps=['go-compile'],
)

k8s_yaml('deployment/test-deployment.yaml')
k8s_resource(
  'nginx',
  objects=[
    'nginx:serviceaccount',
    'workloads:namespace',
  ],
)