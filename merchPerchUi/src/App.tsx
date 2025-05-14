// import { useState } from 'react'
// import { Button } from "@/components/ui/button"
import { Routes, Route } from "react-router-dom"
import ArtistsPage from "./pages/artists"
import ProductsPage from "./pages/products"
import NavBar from "./components/ui/navbar"
import { useRef, useLayoutEffect, useState } from 'react'


function App() {
  const navRef = useRef<HTMLDivElement>(null)
  const [navHeight, setNavHeight] = useState(0)

  useLayoutEffect(() => {
    if (navRef.current) {
      setNavHeight(navRef.current.offsetHeight)
    }
  }, [])

  return (
    <>
      <NavBar ref={navRef} />
      <div style={{ paddingTop: navHeight }}/>
      <Routes>
        <Route path="/" element={<ArtistsPage />} />
        <Route path="/products" element={<ProductsPage />} />
      </Routes>
    </>
  )
}

export default App
