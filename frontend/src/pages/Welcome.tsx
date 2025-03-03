import { useSearchParams, useNavigate } from "@solidjs/router"
import { createEffect, createSignal, Show } from "solid-js"
import toast from "solid-toast"

export default function Welcome() {
  const [searchParams, _] = useSearchParams()
  const navigate = useNavigate()

  createEffect(() => {
    if (searchParams.state === "error") {
      toast.error('Something went wrong\nPlease try again')
      navigate("/")
      return
    }
  })
  let [loading, setLoading] = createSignal(false)
  async function auth() {
    setLoading(true)
    try {
      const response = await fetch(`http://localhost:8080/auth/login`)
      const { url } = (await response.json()) as { url: string }
      window.location.href = url
    } catch (error) {
      console.log(`Error while requesting auth ${error}`)
      toast.error("Something went wrong.\nRetry or contact us")
    }
    setLoading(false)
  }
  return (
    <>
      <main class="h-screen w-full flex gap-2 justify-center items-center">
        <a href="https://github.com/joseph0x45/arcane" target="_blank" class="p-2 w-40 text-center bg-black text-white rounded-md">
          <span>
            Get source Code
          </span>
        </a>
        <button disabled={loading()} onclick={auth} class={`p-2 w-40 bg-black text-white rounded-md ${loading() ? "cursor-not-allowed" : "cursor-pointer"}`}>
          <Show when={loading()}>
            <span class="animate-pulse">
              Login with Github
            </span>
          </Show>
          <Show when={!loading()}>
            <span>
              Login with Github
            </span>
          </Show>
        </button>
      </main>
    </>
  )
}
