export interface FullAccount {
  id: number;
  uuid: number;
  type: string;
  name: string;
  username: string;
  email: string;
  created: string;
}

export interface EditAccount {
  id: number;
  name: string;
  username: string;
  password: string;
  email: string;
}
