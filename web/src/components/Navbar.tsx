import { useAtom } from 'jotai';
import React, { Suspense } from 'react';
import smark from "../assets/smark.png";
import { bookmarksAtom, currentListAtom } from '../state';
import { Bookmark } from '../types';
import Profile from './Profile';

type NavbarProps = {}

const Navbar: React.FC<NavbarProps> = () => {

    const [navList, setNavList] = useAtom(bookmarksAtom);
    const [currentList, setCurrentList] = useAtom(currentListAtom);

    const handleClick = (ele: Bookmark) => {
        console.log('clicked')
        setCurrentList(ele.title);
        // setNavList(curr => curr.map(l => ({ ...l, selected: ele.title == l.title })))
    };

    return (
        <Suspense fallback={null}>
            <div className='h-screen px-4 w-auto max-w-xs flex flex-col'>
                <div className='mt-4 ml-4'>
                    <img src={smark} alt="smark" className='h-10 w-auto' />
                </div>
                <div className='mt-8 w-full h-full max-h-full overflow-y-auto flex flex-col'>
                    {
                        navList.map((ele, i) => (
                            <div
                                onClick={() => handleClick(ele)} key={i}
                                className={`py-2 px-4 ${ele.title == currentList && "bg-dark-gray"} bg-opacity-60 w-full max-w-full hover:bg-dark-gray hover:bg-opacity-70 hover:cursor-pointer rounded-lg font-semibold`}
                            >
                                {ele.title}
                            </div>
                        ))
                    }
                </div>
                <div className='my-4'>
                    <Profile />
                </div>
            </div>
        </Suspense>
    )
}

export default Navbar
