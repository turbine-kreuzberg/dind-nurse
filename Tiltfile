
disable_snapshots()
allow_k8s_contexts(os.getenv("TILT_ALLOW_CONTEXT"))

include('./tests/Tiltfile')

k8s_yaml('deployment/kubernetes.yaml')
k8s_resource(
  'dind',
  port_forwards=['2375', '12375', '40000'],
  labels=["Application"],
  pod_readiness="wait",
)

target='prod'
live_update=[]
if os.environ.get('PROD', '') == '':
  if os.environ.get('DEBUG', '') == '':
    target='build-env'
    live_update=[
      sync('go.mod', '/app/go.mod'),
      sync('go.sum', '/app/go.sum'),
      sync('pkg',    '/app/pkg'),
      sync('main.go', '/app/main.go'),
      run('go install .'),
    ]
  else:
    target='debug-env'
    live_update=[
      sync('go.mod', '/app/go.mod'),
      sync('go.sum', '/app/go.sum'),
      sync('pkg',    '/app/pkg'),
      sync('main.go', '/app/main.go'),
      run('go install -gcflags=\"all=-N -l\" .'),
    ]

docker_build(
  ref='ghcr.io/turbine-kreuzberg/dind-nurse:latest',
  context='.',
  dockerfile='deployment/Dockerfile',
  live_update=live_update,
  target=target,
  only=[ 'go.mod'
       , 'go.sum'
       , 'pkg'
       , 'main.go'
       , 'deployment/entrypoint-delve.sh'
       , 'deployment/entrypoint.sh'
  ],
  ignore=[ '.git'
         , '*/*_test.go'
         , '*.yaml'
  ],
)
