# -*- mode: Python -*-

load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_resource', 'helm_resource', 'helm_repo')
update_settings(max_parallel_updates=10, k8s_upsert_timeout_secs=180)

# Spin up redis
helm_repo('bitnami', 'https://charts.bitnami.com/bitnami', labels='redis')
helm_resource(
  'redis',
  'bitnami/redis',
  port_forwards=["6379:6379"],
  flags=[
    '-f',
    'redis.values.yaml',
  ],
  deps=['redis.values.yaml'],
  resource_deps=['bitnami'],
  labels=['redis'],
)