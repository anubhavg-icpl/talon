'use client'

import { useCallback, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { executeCryptoDecode, type CryptoResult } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'

const OPERATIONS = [
  { group: 'Base Encodings', ops: ['base64_encode', 'base64_decode', 'base32_encode', 'base32_decode', 'base58_encode', 'base58_decode', 'hex_encode', 'hex_decode'] },
  { group: 'Web Encodings', ops: ['url_encode', 'url_decode', 'html_encode', 'html_decode', 'unicode_encode', 'unicode_decode'] },
  { group: 'Classical Ciphers', ops: ['rot13', 'caesar_encode', 'caesar_decode', 'morse_encode', 'morse_decode'] },
  { group: 'Hashing', ops: ['md5_hash', 'sha1_hash', 'sha256_hash', 'sha512_hash'] },
  { group: 'Symmetric Crypto', ops: ['aes_encrypt', 'aes_decrypt', 'des_encrypt', 'des_decrypt'] },
  { group: 'JWT', ops: ['jwt_decode'] },
  { group: 'Auto-Detect', ops: ['auto_decode'] },
]

const NEEDS_KEY = new Set(['aes_encrypt', 'aes_decrypt', 'des_encrypt', 'des_decrypt'])
const NEEDS_SHIFT = new Set(['caesar_encode', 'caesar_decode'])

export default function CryptoToolkit() {
  const [operation, setOperation] = useState('base64_decode')
  const [input, setInput] = useState('')
  const [key, setKey] = useState('')
  const [iv, setIv] = useState('')
  const [shift, setShift] = useState(3)
  const [result, setResult] = useState<CryptoResult | null>(null)
  const [busy, setBusy] = useState(false)

  const handleExecute = useCallback(async () => {
    if (!input.trim()) {
      toast.error('Input is required')
      return
    }
    setBusy(true)
    setResult(null)
    try {
      const res = await executeCryptoDecode({
        operation,
        input,
        key: key || undefined,
        iv: iv || undefined,
        shift: NEEDS_SHIFT.has(operation) ? shift : undefined,
      })
      setResult(res)
      if (!res.success) {
        toast.error(res.error || 'Operation failed')
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Request failed')
    } finally {
      setBusy(false)
    }
  }, [operation, input, key, iv, shift])

  return (
    <div className='grid-bg min-h-screen p-6'>
      <div className='mx-auto max-w-4xl space-y-6'>
        <div>
          <h1 className='font-mono text-2xl font-bold tracking-widest uppercase text-cyan-400'>
            Crypto Toolkit
          </h1>
          <p className='text-muted-foreground mt-1 font-mono text-xs'>
            29 operations: encode, decode, hash, encrypt, decrypt
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>OPERATION</CardTitle>
          </CardHeader>
          <CardContent className='space-y-4'>
            <Select value={operation} onValueChange={(v) => v && setOperation(v)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {OPERATIONS.map((group) => (
                  <div key={group.group}>
                    <div className='micro-label px-2 py-1 text-[10px] text-muted-foreground'>
                      {group.group.toUpperCase()}
                    </div>
                    {group.ops.map((op) => (
                      <SelectItem key={op} value={op}>
                        {op}
                      </SelectItem>
                    ))}
                  </div>
                ))}
              </SelectContent>
            </Select>

            <div className='space-y-2'>
              <Label className='micro-label'>INPUT</Label>
              <Textarea
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder='Enter text to encode/decode/hash/encrypt...'
                className='min-h-[100px] font-mono text-xs'
              />
            </div>

            <div className='grid grid-cols-3 gap-4'>
              {NEEDS_KEY.has(operation) && (
                <div className='space-y-2'>
                  <Label className='micro-label'>KEY</Label>
                  <Input
                    value={key}
                    onChange={(e) => setKey(e.target.value)}
                    placeholder='Encryption key'
                    className='font-mono text-xs'
                  />
                </div>
              )}
              {NEEDS_KEY.has(operation) && operation.startsWith('aes') && (
                <div className='space-y-2'>
                  <Label className='micro-label'>IV (CBC only)</Label>
                  <Input
                    value={iv}
                    onChange={(e) => setIv(e.target.value)}
                    placeholder='16-byte IV'
                    className='font-mono text-xs'
                  />
                </div>
              )}
              {NEEDS_SHIFT.has(operation) && (
                <div className='space-y-2'>
                  <Label className='micro-label'>SHIFT</Label>
                  <Input
                    type='number'
                    value={shift}
                    onChange={(e) => setShift(Number(e.target.value))}
                    className='font-mono text-xs'
                  />
                </div>
              )}
            </div>

            <Button onClick={handleExecute} disabled={busy} className='w-full font-mono tracking-widest uppercase'>
              {busy ? 'EXECUTING...' : 'EXECUTE'}
            </Button>
          </CardContent>
        </Card>

        {result && (
          <Card>
            <CardHeader>
              <CardTitle className='micro-label'>
                RESULT {result.success ? '✓' : '✗'}
              </CardTitle>
            </CardHeader>
            <CardContent>
              {result.success ? (
                <pre className='bg-black/60 max-h-96 overflow-auto rounded-sm p-4 font-mono text-xs text-cyan-400/80 whitespace-pre-wrap'>
                  {result.result}
                </pre>
              ) : (
                <p className='text-destructive font-mono text-xs'>{result.error}</p>
              )}
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  )
}
