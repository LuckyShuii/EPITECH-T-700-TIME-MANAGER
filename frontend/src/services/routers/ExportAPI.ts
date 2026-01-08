import APIHandler from '../APIHandler';

const resource = 'kpi';

export interface KpiExportRequest {
  end_date: string;
  kpi_type: string;
  start_date: string;
  uuid_to_search: string;
}

export interface KpiExportResponse {
  file: string;
  url: string;
}

export default {
  exportKpiData(payload: KpiExportRequest): Promise<KpiExportResponse> {
    return APIHandler.post(`${resource}/export`, payload);
  }
};