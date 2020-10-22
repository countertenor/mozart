export const NOTIFICATION_INIT = {
    title: '',
    subtitle: '',
    caption: '',
    kind: '',
    isVisible: false
  };
  
  export const IBM_CONFIG_INIT = {
    clusterName: '',
    regionItem: null,
    zone: '',
    zoneDropdownItems: [],
    apiKey: '',
    resourceGroup: '',
    privateVLAN: '',
    publicVLAN: '',
  };
  
  export const AWS_CONFIG_INIT = {
    clusterName: '',
    regionItem: null,
    accessKey: '',
    accessKeyId: '',
    dns: '',
    rhPullSecret: '',
    ocpWorkerSingleZoneItem: { label: 'True', value: 'true' },
    clusterTypeItem: { label: 'Public', value: 'PUBLIC' },
    subnetsList: ''
  };
  
  export const DEPLOYMENT_CONFIG_INIT = {
    type: 'production',
    nodes: 1,
    selectedNodes: [],
    organization: '',
    instances: 1,
    databases: 1,
    replicated: false,
    clusterName: '',
    databaseName: '',
    instanceName: '',
    instanceSecret: '',
    fencedName: '',
    fencedSecret: '',
    virtualIP: '',
    cidrNetmask: 0,
    nic: ''
    // os: 'redhat 7.8',
  };
  
  export function CONFIG_REDUCER(state, action) {
    return {
      ...state,
      [action.type]: action.value
    };
  }
  
  export function NOTIFICATION_REDUCER(state, action) {
    switch (action.type) {
    case 'setNotification':
      return {
        ...action.notification,
        isVisible: true // if you want to hide, use clear
      };
    case 'clearNotification':
      return NOTIFICATION_INIT;
    default: // Do nothing on unknwwn
      return state;
    }
  }
  