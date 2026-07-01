import { Input } from "~/components/auth/input";
import { Footer } from "~/components/auth/footer";
import { useEffect, useRef, useState } from "react";
import type { Route } from "./+types/home";
import { reqServer } from "~/lib/apiFetch";
import type { ApiRequest } from "~/lib/apiFetch";
import { Header } from "~/components/auth/header";
export function meta({}: Route.MetaArgs) {
  return [{ title: "Tano" }, { name: "", content: "Welcome to TANO!" }];
}

export default function Register() {
  const [selected, setSelected] = useState(0);
  const [reveal, setReveal] = useState(false);

  function shortcutHandler(e: React.KeyboardEvent) {
    if (e.key === "ArrowDown") {
      setSelected(Math.min(selected + 1, 1));
    }
    if (e.key === "ArrowUp") {
      setSelected(Math.max(selected - 1, 0));
    }
    if (e.key === "Enter") {
      console.log("submit");
      submitHandler({
        email: emailRef.current?.value ?? "",
        password: passwordRef.current?.value ?? "",
      });
    }
    if (e.key === "F2") {
      const nextReveal = !reveal;
      setReveal(nextReveal);
    }
  }

  const emailRef = useRef<HTMLInputElement>(null);
  const passwordRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    switch (selected) {
      case 0:
        emailRef.current?.focus();
        break;
      case 1:
        passwordRef.current?.focus();
        break;
    }
  }, [selected]);

  const [err, setErr] = useState(false);
  const [errString, setErrString] = useState("");

  async function submitHandler(data: { email: string; password: string }) {
    if (data.email.trim() === "") {
      setErr(true);
      setErrString("empty email field");
      return;
    }
    if (data.password.trim() === "") {
      setErr(true);
      setErrString("empty password field");
      return;
    }
    const emailErr = emailValidation(data.email);
    if (emailErr != "") {
      setErr(true);
      setErrString(emailErr);
      setTimeout(() => {
        setErr(true);
        setErrString(res.err);
      }, 3000);
      return;
    }
    const passwordErr = passwordValidation(data.password);
    if (passwordErr != "") {
      setErr(true);
      setErrString(passwordErr);
      setTimeout(() => {
        setErr(true);
        setErrString(res.err);
      }, 3000);
      return;
    }
    let T: ApiRequest = {
      body: {
        email: data.email,
        password: data.password,
      },
      method: "POST",
      route: "user/register",
    };
    let res = await reqServer(T);
    if (!res.ok) {
      console.log(res.err);
      setErr(true);
      setErrString(res.err);
      setTimeout(() => {
        setErr(true);
        setErrString(res.err);
      }, 3000);
      return;
    }
    setErr(false);
    setErrString("");
    return;
  }

  return (
    <div
      className="p-16 font-mono flex-col justify-center content-center h-screen bg-gunmetal-950 relative"
      onKeyDown={shortcutHandler}
      tabIndex={0}
    >
      <Header name={"register_"} />
      <form className="flex-col p-8" onKeyDown={shortcutHandler}>
        <Input labelName="enter your email" inputType="text" ref={emailRef} />
        <Input
          labelName="enter your password"
          inputType={reveal ? "text" : "password"}
          ref={passwordRef}
        />
      </form>
      <Footer err={err} errValue={errString} />
    </div>
  );
}

function passwordValidation(password: string): string {
  if (password.length < 8) {
    return "password should have a minimum of 8 characters";
  }
  let uppercase = new RegExp("[A-Z]");
  let lowercase = new RegExp("[a-z]");
  let numbers = new RegExp("[0-9]");
  const special = new RegExp("[^A-Za-z0-9]");
  if (!uppercase.test(password)) {
    return "password must have atleast one uppercase character";
  }
  if (!lowercase.test(password)) {
    return "password must have atleast one lowercase character";
  }
  if (!numbers.test(password)) {
    return "password must have atleast one digit";
  }
  if (!special.test(password)) {
    return "password must have atleast one special character";
  }
  return "";
}

function emailValidation(email: string): string {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  if (!emailRegex.test(email)) {
    return "Invalid email address.";
  }

  return "";
}
