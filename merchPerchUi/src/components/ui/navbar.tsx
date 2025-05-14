// components/NavBar.tsx
import React from 'react'
import { Link, useLocation } from 'react-router-dom'

const NavBar = React.forwardRef<HTMLDivElement>((_, ref) => (
    <nav ref={ref} className="fixed top-0 left-0 w-full z-50 bg-white shadow px-4 py-2">
    <div className="relative flex items-center">
      {/* Left: Navigation links */}
      <div className="flex gap-4">
        <Link to="/" className="text-sm font-medium">Artists</Link>
        <Link to="/products" className="text-sm font-medium">Products</Link>
      </div>

      {/* Center: Merch Perch */}
      <div className="absolute left-1/2 -translate-x-1/2 text-lg font-bold">
        Merch Perch
      </div>
    </div>
    </nav>
))

export default NavBar
