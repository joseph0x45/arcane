import { Router, Route } from "@solidjs/router"
import Welcome from "./pages/Welcome"
import { Toaster } from "solid-toast"
import Home from "./pages/Home"

export default function App() {
  return (
    <>
      <Toaster position="bottom-right" gutter={8} />
      <Router>
        <Route path="/" component={Welcome} />
        <Route path="/home" component={Home} />
      </Router>
    </>
  )
}
