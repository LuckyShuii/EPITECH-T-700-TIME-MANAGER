export interface Employee {
  uuid: string
  first_name: string
  last_name: string
  email: string
  username: string
  phone_number: string
  roles: string[]
  weekly_hours: number  // 35, 39, etc.
  status: 'active' | 'inactive'
}

export interface EmployeeUpdateData {
  email: string
  phone_number: string
  roles: string[]
  weekly_hours: number
  status: 'active' | 'inactive'
}