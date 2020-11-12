// '' means notstarted, 'error', 'success', 'running', 'skipped' (skipped based on steps failing)

export const STATUS_DUMMY = Object.freeze({ // eslint-disable-line
    steps: [
      {
        module: 'generated/install/10-linbit',
        tasks: [
          {
            taskName: '00-set-env.sh',
            status: {
              lastSuccessTime: '2020-09-20 18:38:33.163370832 -0700 PDT m=+1.300127619',
              lastErrorTime: '',
              state: 'success',
              logFilePath: '/var/log/db2-orchestrator/2020-09-20--18-38-31.888-00-set-env.log'
            }
          },
          {
            taskName: '10-create-linstor-conf.sh',
            status: {
              lastSuccessTime: '2020-09-20 18:39:23.115400019 -0700 PDT m=+51.252156853',
              lastErrorTime: '',
              state: 'success',
              logFilePath: '/var/log/db2-orchestrator/2020-09-20--18-39-03.585-10-create-linstor-conf.log'
            }
          }
        ]
      },
      {
        module: 'generated/install/20-pacemaker',
        tasks: [
          {
            taskName: '00-set-env.sh',
            status: {
              lastSuccessTime: '2020-09-20 18:38:33.163370832 -0700 PDT m=+1.300127619',
              lastErrorTime: '',
              state: 'success',
              logFilePath: '/var/log/db2-orchestrator/2020-09-20--18-38-31.888-00-set-env.log'
            }
          },
          {
            taskName: '10-create-linstor-conf.sh',
            status: {
              lastSuccessTime: '2020-09-20 18:39:23.115400019 -0700 PDT m=+51.252156853',
              lastErrorTime: '',
              state: 'success',
              logFilePath: '/var/log/db2-orchestrator/2020-09-20--18-39-03.585-10-create-linstor-conf.log'
            }
          }
        ]
      },
      {
        module: 'generated/install/30-db2',
        tasks: [
          {
            taskName: '00-set-env.sh',
            status: {
              lastSuccessTime: '2020-09-20 18:38:33.163370832 -0700 PDT m=+1.300127619',
              lastErrorTime: '',
              state: 'success',
              logFilePath: '/var/log/db2-orchestrator/2020-09-20--18-38-31.888-00-set-env.log'
            }
          },
          {
            taskName: '10-create-linstor-conf.sh',
            status: {
              lastSuccessTime: '2020-09-20 18:39:23.115400019 -0700 PDT m=+51.252156853',
              lastErrorTime: '',
              state: 'success',
              logFilePath: '/var/log/db2-orchestrator/2020-09-20--18-39-03.585-10-create-linstor-conf.log'
            }
          }
        ]
      }
    ]
  });
  