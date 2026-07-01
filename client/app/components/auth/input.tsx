export function Input(values: {
  labelName: string, inputType: string,
  ref?: React.RefObject<HTMLInputElement | null>
}) {
  return (
    <div className="m-4">
      <input type={values.inputType} placeholder={values.labelName} ref={values.ref}
        className="w-full text-platinum-500 focus:text-sage-green-500 border border-gunmetal-800 rounded-2xl p-4 " />
    </div>
  );
}
