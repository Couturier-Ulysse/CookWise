import Hat from '../assets/Hat.png';

export default function Header() {
    return (
        <header className='mt-0' style={{ top: '16px', backgroundColor: 'transparent' }} >
            <div className="flex flex-col items-center max-w-5xl mx-auto relative">
                <img
                    src={Hat}
                    alt="Stove"
                    className="opacity-70 w-32 md:w-48"
                    style={{ marginTop: '75px' }}
                />
                <h1 className="text-5xl font-semibold uppercase text-amber-800 font-title tracking-widest relative z-10">
                    COOKWISE
                </h1>
                <span className="text-lg italic text-gray-500 mt-1 relative z-10">
                    ' Your weekly meal planner '
                </span>
            </div>
        </header>
    );
}
