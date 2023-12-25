import { FileWithPath } from "@mantine/dropzone";

export type FileStoreType = {
  files: FileWithPath[];
  setFiles: (file: FileWithPath[]) => void;
  clearFiles: () => void;
};
