import type { Metadata } from "next";
import { DocsPage } from "@/components/docs-page";

export const metadata: Metadata = {
  title: "Hermes docs",
  description: "Install Hermes, hand off work between coding agents, pick a session, resume, and every flag.",
};

export default function Page() {
  return <DocsPage />;
}
