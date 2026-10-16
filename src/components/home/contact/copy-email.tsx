// SPDX-License-Identifier: MPL-2.0
"use client";

import { useEffect, useState } from "react";
import { buttonClass } from "@/components/ui/button";

type CopyEmailProps = {
  email: string;
  /** The element showing the address, selected for the visitor when the browser refuses clipboard access. */
  addressId: string;
  labels: { copy: string; copied: string; selected: string };
};

/** Copies the address instead of relying on a mail app, which many computers (and expo screens) do not have. */
export function CopyEmail({ email, addressId, labels }: CopyEmailProps) {
  const [state, setState] = useState<"copy" | "copied" | "selected">("copy");

  useEffect(() => {
    if (state === "copy") return undefined;
    const timer = setTimeout(() => setState("copy"), 3000);
    return () => clearTimeout(timer);
  }, [state]);

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(email);
      setState("copied");
    } catch {
      const address = document.getElementById(addressId);
      if (address) window.getSelection()?.selectAllChildren(address);
      setState("selected");
    }
  };

  return (
    <button type="button" onClick={copy} className={buttonClass()}>
      <span aria-live="polite">{labels[state]}</span>
    </button>
  );
}
