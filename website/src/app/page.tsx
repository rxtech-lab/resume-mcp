import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Terminal,
  FileText,
  Database,
  Code,
  Palette,
  Wrench,
} from "lucide-react";
import { DownloadSection } from "@/components/download-section";

export default function Home() {
  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <header className="pt-16 pb-24">
          <div className="text-center">
            <h1 className="text-4xl sm:text-5xl font-bold text-foreground mb-4">
              Resume MCP
            </h1>
            <p className="text-lg sm:text-xl text-muted-foreground max-w-2xl mx-auto">
              A Model Context Protocol server for AI agents to manage resume
              data and generate PDF previews
            </p>
          </div>
        </header>

        <DownloadSection />

        <section className="mb-24">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-foreground mb-4">
              Features
            </h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              Comprehensive MCP tools for resume management and PDF generation
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            <Card>
              <CardHeader>
                <Terminal className="w-8 h-8 mb-2 text-muted-foreground" />
                <CardTitle className="text-lg">Resume Management</CardTitle>
                <CardDescription>
                  Create, update, and manage resume data with flexible JSON
                  features
                </CardDescription>
              </CardHeader>
            </Card>

            <Card>
              <CardHeader>
                <FileText className="w-8 h-8 mb-2 text-muted-foreground" />
                <CardTitle className="text-lg">PDF Generation</CardTitle>
                <CardDescription>
                  Generate beautiful PDF previews with customizable templates
                </CardDescription>
              </CardHeader>
            </Card>

            <Card>
              <CardHeader>
                <Code className="w-8 h-8 mb-2 text-muted-foreground" />
                <CardTitle className="text-lg">MCP Tools</CardTitle>
                <CardDescription>
                  Comprehensive set of tools for AI agents to interact with
                  resume data
                </CardDescription>
              </CardHeader>
            </Card>

            <Card>
              <CardHeader>
                <Database className="w-8 h-8 mb-2 text-muted-foreground" />
                <CardTitle className="text-lg">Local Storage</CardTitle>
                <CardDescription>
                  SQLite database for secure local data storage with automatic
                  migrations
                </CardDescription>
              </CardHeader>
            </Card>

            <Card>
              <CardHeader>
                <Wrench className="w-8 h-8 mb-2 text-muted-foreground" />
                <CardTitle className="text-lg">REST API</CardTitle>
                <CardDescription>
                  Built-in HTTP server for serving HTML previews and template
                  rendering
                </CardDescription>
              </CardHeader>
            </Card>

            <Card>
              <CardHeader>
                <Palette className="w-8 h-8 mb-2 text-muted-foreground" />
                <CardTitle className="text-lg">Templates</CardTitle>
                <CardDescription>
                  Flexible Go template system with CSS styling support
                </CardDescription>
              </CardHeader>
            </Card>
          </div>
        </section>

        <section className="mb-24">
          <Card>
            <CardHeader className="text-center">
              <CardTitle className="text-2xl">Quick Start</CardTitle>
              <CardDescription>
                Get started with Resume MCP in three simple steps
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid md:grid-cols-3 gap-8 max-w-4xl mx-auto">
                <div className="text-center">
                  <div className="w-12 h-12 bg-muted rounded-full flex items-center justify-center mx-auto mb-4">
                    <span className="text-lg font-semibold">1</span>
                  </div>
                  <h3 className="text-lg font-semibold mb-2">Install</h3>
                  <p className="text-muted-foreground">
                    Download and install the macOS package
                  </p>
                </div>

                <div className="text-center">
                  <div className="w-12 h-12 bg-muted rounded-full flex items-center justify-center mx-auto mb-4">
                    <span className="text-lg font-semibold">2</span>
                  </div>
                  <h3 className="text-lg font-semibold mb-2">Configure</h3>
                  <p className="text-muted-foreground">
                    Set up MCP server in your AI agent
                  </p>
                </div>

                <div className="text-center">
                  <div className="w-12 h-12 bg-muted rounded-full flex items-center justify-center mx-auto mb-4">
                    <span className="text-lg font-semibold">3</span>
                  </div>
                  <h3 className="text-lg font-semibold mb-2">Build</h3>
                  <p className="text-muted-foreground">
                    Start creating and managing resumes
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </section>

        <footer className="pb-12 text-center">
          <div className="flex justify-center space-x-6 mb-4">
            <a
              href="https://github.com/rxtech-lab/resume-mcp"
              className="text-muted-foreground hover:text-foreground transition-colors"
              target="_blank"
              rel="noopener noreferrer"
            >
              GitHub
            </a>
            <a
              href="https://github.com/rxtech-lab/resume-mcp/releases"
              className="text-muted-foreground hover:text-foreground transition-colors"
              target="_blank"
              rel="noopener noreferrer"
            >
              Releases
            </a>
            <a
              href="https://github.com/rxtech-lab/resume-mcp/issues"
              className="text-muted-foreground hover:text-foreground transition-colors"
              target="_blank"
              rel="noopener noreferrer"
            >
              Issues
            </a>
          </div>
          <p className="text-muted-foreground text-sm">Built by RxTech Lab</p>
        </footer>
      </div>
    </div>
  );
}
