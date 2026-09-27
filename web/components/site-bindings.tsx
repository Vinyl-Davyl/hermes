"use client";

import { useEffect } from "react";
import { bindSite } from "@/lib/runtime";

export function SiteBindings() {
  useEffect(() => bindSite(), []);
  return null;
}
