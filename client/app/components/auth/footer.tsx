export function Footer(val: { err: boolean; errValue: string }) {
  if (!val.err) {
    return (
      <div className="text-center text-platinum-500 text-xs bottom-0 left-0 absolute px-8 py-4 w-full bg-gunmetal-900">
        <kbd className="text-sage-green-500">↑/↓</kbd> to move between fields *
        <kbd className="text-sage-green-500"> enter</kbd> to submit *
        <kbd className="text-sage-green-500"> F2</kbd> to view password
      </div>
    );
  }
  return (
    <div className="text-center text-platinum-200 text-xs bottom-0 left-0 absolute px-8 py-4 w-full bg-classic-crimson-800">
      <p> {val.errValue}</p>
    </div>
  );
}
