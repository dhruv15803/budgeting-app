import { Component, type ReactNode } from "react"
import { AlertTriangleIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

interface State {
  error: Error | null
}

interface Props {
  children: ReactNode
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null }

  static getDerivedStateFromError(error: Error) {
    return { error }
  }

  componentDidCatch(error: Error) {
    if (import.meta.env.DEV) console.error("ErrorBoundary:", error)
  }

  reset = () => this.setState({ error: null })

  render() {
    if (this.state.error) {
      return (
        <div className="flex min-h-screen items-center justify-center p-6">
          <Card className="w-full max-w-md">
            <CardHeader>
              <div className="flex items-center gap-2">
                <div className="bg-destructive/10 text-destructive flex size-9 items-center justify-center rounded-md">
                  <AlertTriangleIcon className="size-4" />
                </div>
                <CardTitle>Something broke</CardTitle>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <p className="text-muted-foreground text-sm">{this.state.error.message}</p>
              <div className="flex gap-2">
                <Button onClick={this.reset} variant="outline">
                  Try again
                </Button>
                <Button onClick={() => (window.location.href = "/")}>Go home</Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )
    }
    return this.props.children
  }
}
