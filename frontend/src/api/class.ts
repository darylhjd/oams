export type ClassesGetResponse = {
  classes: Class[];
};

export type ClassGetResponse = {
  class: Class;
};

export type Class = {
  id: number;
  code: string;
  year: number;
  semester: string;
  programme: string;
  au: number;
  created_at: Date;
  updated_at: Date;
};
