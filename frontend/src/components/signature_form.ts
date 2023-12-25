export function signatureInputForm() {
  return {
    initialValues: {
      signature: "",
    },
    validate: {
      signature: (value: string) =>
        value.length == 0 ? "Signature cannot be empty" : null,
    },
  };
}
