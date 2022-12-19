import React, {useRef, useState} from "react"
import {AiOutlineSearch} from "react-icons/ai"

interface ThreadSearchProps {
  setSearchTerm: React.Dispatch<React.SetStateAction<string>>
}

export default function ThreadSearch({setSearchTerm}: ThreadSearchProps) {
  const searchRef = useRef<HTMLInputElement>()
  const [toggled, setToggled] = useState(false)

  const toggledStyles = `p-2 border border-primary rounded-br scale-x-100`
  const untoggledStyles = `scale-x-0 origin-right`

  function onChange() {
    const searchTerm = searchRef.current.value.trim()
    setSearchTerm(searchTerm)
  }

  return (
    <div className="flex items-center phone:flex-row-reverse">
      <input
        type="text"
        id="search"
        className={`${toggled ? toggledStyles : untoggledStyles} transition-all duration-500 phone:w-[100px]`}
        onBlur={() => setToggled(false)}
        ref={searchRef}
        placeholder="Search by title"
        onChange={onChange}/>
      <AiOutlineSearch size={30} className={`text-primary hover:cursor-pointer ${toggled ? 'hidden' : ''}`} onClick={() => {
        document.getElementById("search").focus()
        setToggled(!toggled)
      }}/>
    </div>
  )
}
