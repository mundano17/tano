import { useNavigate } from "react-router";

export function MenuItem(value: {
  name: string;
  shortcut: string;
  route: string;
}) {
  let r = useNavigate();
  return (
    <button
      className="flex space-x-4 text-platinum-400 hover:text-classic-crimson-500 hover:"
      onClick={() => r(value.route)}
    >
      <pre className="text-">{value.name}</pre>
      <kbd>{value.shortcut}</kbd>
    </button>
  );
}
