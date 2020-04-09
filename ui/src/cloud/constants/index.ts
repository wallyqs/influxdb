import {DashboardTemplate} from 'src/types'

export const RATE_LIMIT_ERROR_STATUS = 429

export const RATE_LIMIT_ERROR_TEXT =
  'Oops. It looks like you have exceeded the query limits allowed as part of your plan. If you would like to increase your query limits, reach out to support@influxdata.com.'

export const ASSET_LIMIT_ERROR_STATUS = 403

export const ASSET_LIMIT_ERROR_TEXT =
  'Oops. It looks like you have exceeded the asset limits allowed as part of your plan. If you would like to increase your limits, reach out to support@influxdata.com.'

const WebsiteMonitoringDashboardTemplate = async (name: string) => {
  const websiteMonitoringTemplate = (await import('src/cloud/constants/websiteMonitoringTemplate')) as any
  websiteMonitoringTemplate.content.data.attributes.name = name
  return websiteMonitoringTemplate as DashboardTemplate
}

export const WebsiteMonitoringBucket = 'Website Monitoring Bucket'

export const DemoDataDashboards = {
  [WebsiteMonitoringBucket]: 'Website Monitoring Demo Data Dashboard',
}

export const DemoDataTemplates = {
  [WebsiteMonitoringBucket]: WebsiteMonitoringDashboardTemplate(
    DemoDataDashboards[WebsiteMonitoringBucket]
  ),
}
