import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Download, Github } from "lucide-react";

interface GitHubRelease {
  tag_name: string;
  assets: Array<{
    name: string;
    browser_download_url: string;
  }>;
}

async function getLatestRelease(): Promise<GitHubRelease | null> {
  try {
    const response = await fetch(
      "https://api.github.com/repos/rxtech-lab/resume-mcp/releases/latest",
      {
        next: { revalidate: 3600 }, // Cache for 1 hour
      }
    );

    if (!response.ok) {
      throw new Error("Failed to fetch release data");
    }

    return await response.json();
  } catch (error) {
    console.error("Error fetching release:", error);
    return null;
  }
}

function getMacOSDownloadUrl(release: GitHubRelease | null): string | null {
  if (!release) return null;

  const macOSAsset = release.assets.find(
    (asset) =>
      asset.name.includes("macOS") &&
      asset.name.includes("arm64") &&
      asset.name.endsWith(".pkg")
  );

  return macOSAsset?.browser_download_url || null;
}

export async function DownloadSection() {
  const release = await getLatestRelease();
  const downloadUrl = getMacOSDownloadUrl(release);

  return (
    <section className="mb-24">
      <Card className="text-center">
        <CardHeader>
          <CardTitle className="text-2xl">Download for macOS</CardTitle>
          <CardDescription>
            Get the latest version and start building resume management tools
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button asChild size="lg" disabled={!downloadUrl}>
              {downloadUrl ? (
                <a href={downloadUrl} className="gap-2">
                  <Download className="w-4 h-4" />
                  Download {release?.tag_name || "Latest"}
                </a>
              ) : (
                <div className="gap-2">
                  <Download className="w-4 h-4" />
                  Download Unavailable
                </div>
              )}
            </Button>
            <Button asChild variant="outline" size="lg">
              <a
                href="https://github.com/rxtech-lab/resume-mcp"
                target="_blank"
                rel="noopener noreferrer"
                className="gap-2"
              >
                <Github className="w-4 h-4" />
                View on GitHub
              </a>
            </Button>
          </div>
          {!downloadUrl && (
            <p className="text-sm text-muted-foreground mt-2">
              Unable to fetch latest release. Please visit GitHub for manual
              download.
            </p>
          )}
        </CardContent>
      </Card>
    </section>
  );
}
