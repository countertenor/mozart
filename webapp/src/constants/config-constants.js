export const AWS_REGIONS = [
    { label: 'US East (Ohio)', value: 'us-east-2' },
    { label: 'US East (N. Virginia)', value: 'us-east-1' },
    { label: 'US West (N. California)', value: 'us-west-1' },
    { label: 'US West (Oregon)', value: 'us-west-2' },
    { label: 'Africa (Cape Town)', value: 'af-south-1' },
    { label: 'Asia Pacific (Hong Kong)', value: 'ap-east-1' },
    { label: 'Asia Pacific (Mumbai)', value: 'ap-south-1' },
    { label: 'Asia Pacific (Osaka-Local)', value: 'ap-northeast-3' },
    { label: 'Asia Pacific (Seoul)', value: 'ap-northeast-2' },
    { label: 'Asia Pacific (Singapore)', value: 'ap-southeast-1' },
    { label: 'Asia Pacific (Sydney)', value: 'ap-southeast-2' },
    { label: 'Asia Pacific (Tokyo)', value: 'ap-northeast-1' },
    { label: 'Canada (Central)', value: 'ca-central-1' },
    { label: 'China (Beijing)', value: 'cn-north-1' },
    { label: 'China (Ningxia)', value: 'cn-northwest-1' },
    { label: 'Europe (Frankfurt)', value: 'eu-central-1' },
    { label: 'Europe (Ireland)', value: 'eu-west-1' },
    { label: 'Europe (London)', value: 'eu-west-2' },
    { label: 'Europe (Milan)', value: 'eu-south-1' },
    { label: 'Europe (Paris)', value: 'eu-west-3' },
    { label: 'Europe (Stockholm)', value: 'eu-north-1' },
    { label: 'Middle East (Bahrain)', value: 'me-south-1' },
    { label: 'South America (São Paulo)', value: 'sa-east-1' },
    { label: 'AWS GovCloud (US-East)', value: 'us-gov-east-1' },
    { label: 'AWS GovCloud (US)', value: 'us-gov-west-1' }
  ];
  
  // Unused for now, think this is for AWS so might need later so don't remove
  // export const ZONES = [
  //   'ams03',
  //   'che01',
  //   'hkg02',
  //   'mel01',
  //   'mex01',
  //   'mil01',
  //   'mon01',
  //   'osl01',
  //   'par01',
  //   'sao01',
  //   'seo01',
  //   'sjc03',
  //   'sjc04',
  //   'sng01',
  //   'tor01',
  //   'dal10',
  //   'dal12',
  //   'dal13',
  //   'fra02',
  //   'fra04',
  //   'fra05',
  //   'lon02',
  //   'lon04',
  //   'lon05',
  //   'lon06',
  //   'syd01',
  //   'syd04',
  //   'syd05',
  //   'tok02',
  //   'tok04',
  //   'tok05',
  //   'wdc04',
  //   'wdc06',
  //   'wdc',
  // ]; // TODO: Organize for input
  
  export const IBM_REGIONS = [
    { label: 'Asia Pacific', value: 'Asia Pacific' },
    { label: 'Europe', value: 'Europe' },
    { label: 'North America', value: 'North America' },
    { label: 'South America', value: 'South America' },
  ];
  
  export const IBM_ZONES = {
    'Asia Pacific': ['che01', 'hkg02', 'mel01', 'seo01', 'sng01', 'syd01', 'syd04', 'syd05', 'tok02', 'tok04', 'tok05'],
    Europe: ['ams03', 'fra02', 'fra04', 'fra05', 'lon02', 'lon04', 'lon05', 'lon06', 'mil01', 'osl01', 'par01'],
    'North America': ['dal10', 'dal12', 'dal13', 'hou02', 'mex01',
      'mon01', 'sjc03', 'sjc04', 'tor01', 'wdc04', 'wdc06', 'wdc07'],
    'South America': ['sao01']
  };
  
  export const PROVIDERS = ['IBM Cloud', 'Amazon Web Services'];
  