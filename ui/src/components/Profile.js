import React, { useEffect, useLayoutEffect, useState } from 'react'
import { GetProfileData } from '../utils'

const Profile = ({toLogout}) => {
  const [data, setData]=useState()

  useLayoutEffect(()=>{
    GetProfileData(setData)
  },[])
  console.log(data)
  return (
    <div>
        <div className="details-box">
          <p className=''>Message:{data?.message}</p>
          <p>EMAIL: <span>{data?.data.email}</span></p>
          <p>USERNAME: <span>{data?.data.username}</span></p>
          <p>ROLE: <span>{data?.data.role}</span></p>
          <p>ROLE: <span>{data?.data.userid}</span></p>
          <button className="logout" onClick={toLogout}>Logout</button> 
        </div>
    </div>
  )
}

export default Profile