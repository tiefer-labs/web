// SPDX-License-Identifier: MPL-2.0
import type { ReactNode } from "react";
import { cn } from "@/lib/cn";

export const inputClass =
  "w-full rounded-[2px] border border-ink-60 bg-paper px-3.5 py-3 font-text text-base leading-[1.4] text-ink transition-colors duration-150 hover:border-ink aria-invalid:border-2 aria-invalid:border-ink aria-invalid:px-[13px] aria-invalid:py-[11px]";

/** The `id`, `aria-invalid` and `aria-describedby` a control needs to point at its hint and error. */
export function describe(id: string, { error, hint }: { error?: string; hint?: boolean }) {
  const describedBy = [hint && `${id}-hint`, error && `${id}-err`].filter(Boolean).join(" ");
  return { id, "aria-invalid": error ? true : undefined, "aria-describedby": describedBy || undefined };
}

type FieldProps = {
  id: string;
  label: ReactNode;
  optional?: string;
  hint?: string;
  error?: string;
  errorPrefix: string;
  className?: string;
  children: ReactNode;
};

export function Field({ id, label, optional, hint, error, errorPrefix, className, children }: FieldProps) {
  return (
    <div className={cn("grid gap-2", className)}>
      <label htmlFor={id} className="text-small font-medium text-ink">
        {label}
        {optional && <span className="font-normal text-ink-60"> ({optional})</span>}
      </label>
      {children}
      {hint && (
        <p id={`${id}-hint`} className="text-[0.875rem] text-ink-60">
          {hint}
        </p>
      )}
      <FieldError id={id} error={error} prefix={errorPrefix} />
    </div>
  );
}

export function FieldError({
  id,
  error,
  prefix,
  className,
}: {
  id: string;
  error?: string;
  prefix: string;
  className?: string;
}) {
  if (!error) return null;
  return (
    <p id={`${id}-err`} className={cn("text-[0.875rem] font-medium text-ink", className)}>
      <span className="sr-only">{prefix}</span> {error}
    </p>
  );
}
