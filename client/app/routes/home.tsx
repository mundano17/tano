import { useNavigate } from "react-router";
import { useEffect, useRef } from "react";
import type { Route } from "./+types/home";
import { MenuItem } from "~/components/home/menuItem";
export function meta({}: Route.MetaArgs) {
  return [{ title: "Tano" }, { name: "", content: "Welcome to React Router!" }];
}

export default function Home() {
  let logo = `
████████╗ █████╗ ███╗   ██╗ ██████╗
╚══██╔══╝██╔══██╗████╗  ██║██╔═══██╗
   ██║   ███████║██╔██╗ ██║██║   ██║
   ██║   ██╔══██║██║╚██╗██║██║   ██║
   ██║   ██║  ██║██║ ╚████║╚██████╔╝
   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝
`;
  const route = useNavigate();
  const divRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    divRef.current?.focus();
  }, []);
  return (
    <div
      ref={divRef}
      className="h-screen flex-col content-center justify-center gap-4 bg-gunmetal-950 font-mono"
      onKeyDown={(e) => {
        if (e.key.toLowerCase() == "l") {
          route("/login");
        }
        if (e.key.toLowerCase() == "r") {
          route("/register");
        }
      }}
      tabIndex={0}
    >
      <div className="text-center py-8 text-classic-crimson-500">
        <pre>{logo}</pre>
      </div>
      <nav className="grid-rows-2 gap-2 grid justify-center">
        <MenuItem name={"Login   "} shortcut={"[L]"} route="/login" />
        <MenuItem name="Register" shortcut="[R]" route="/register" />
      </nav>
    </div>
  );
}
