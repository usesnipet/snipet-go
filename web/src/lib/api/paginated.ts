export type Paginated<T> = {
  data: T[];
  total: number;
  take: number;
  skip: number;
}