export interface Employee {
  uuid: string
  first_name: string
  last_name: string
  email: string
  username: string
  phone_number: string
  roles: string[]
  weekly_rate: number
  weekly_rate_name: string
  status: 'active' | 'inactive'
  work_session_status?: string
  weekly_rate_uuid?: string  // Ajout optionnel pour éviter les erreurs
}

export interface EmployeeUpdateData {
  first_name: string
  last_name: string
  email: string
  phone_number: string
  roles: string[]
  weekly_rate_uuid: string     
  status: 'active' | 'inactive'
  username: string             
}