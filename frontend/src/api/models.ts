export type LoginResponse = {
  redirect_url: string;
}

export type User = {
  id: string;
  name: string;
  email: string;
  role: string;
  createdAt: Date;
  updatedAt: Date;
}
