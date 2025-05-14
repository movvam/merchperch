// import { useState } from 'react'
// import { Button } from "@/components/ui/button"
import ArtistJsonDataDisplay from "@/components/artistJsonTable"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Routes, Route } from "react-router-dom"
import ArtistsPage from "./pages/artists"
import ProductsPage from "./pages/products"
import NavBar from "./components/ui/navbar"


function App() {
  return (
    <>
      <NavBar />
      <Routes>
        <Route path="/" element={<ArtistsPage />} />
        <Route path="/products" element={<ProductsPage />} />
      </Routes>
    </>
  )
}

export default App
