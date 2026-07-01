export function Header(val: { name: string }) {
  return (
    <div className="text-center text-platinum-500 text-xl top-0 left-0 absolute px-8 py-4 w-full bg-gunmetal-900">
      {val.name}
    </div>
  );
}
