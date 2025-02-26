/* @refresh reload */
import { render } from 'solid-js/web'
import './index.css'
import App from './App.tsx'
import { Router, Route } from "@solidjs/router"

const wrapper = document.getElementById("root")

if (!wrapper) {
  throw new Error("Wrapper div not found")
}

render(() =>
  <Router>
    <Route path="/" component={App} />
  </Router>, wrapper)
