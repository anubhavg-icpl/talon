const Footer = () => {
  return (
    <footer className='op-footer'>
      <div className='text-muted-foreground mx-auto flex size-full max-w-360 items-center justify-between gap-3 px-4 py-3 font-mono text-[10px] tracking-widest uppercase max-sm:flex-col sm:gap-6 sm:px-6'>
        <p>
          TALON <span className='text-primary'>//</span> OPERATOR — RECON · EXPLOIT · JUDGE
        </p>
        <p className='text-primary/70 max-sm:hidden'>AUTHORIZED TARGETS ONLY</p>
      </div>
    </footer>
  )
}

export default Footer
