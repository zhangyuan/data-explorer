export interface QueryParams {
  [key: string]: string;
}

export interface QueryResult {
  column_names: string[]
  column_types: string[]
  records: any[][]
}

export interface SqlQuery {
  id: number
  duration: number
  query: string
  sql: string
  result: QueryResult
  params: QueryParams
}

export interface IssueItem {
  id: number
  title: string
  description: string
  created_at: string
  updated_at: string
}

export interface ItemSection {
  id: number
  header: string
  body: string
  footer: string
  updated_at: string
  queries: SqlQuery[] | undefined
}
