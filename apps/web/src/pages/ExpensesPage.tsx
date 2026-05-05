import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts'
import { ArrowLeft, Plus, Trash2 } from 'lucide-react'
import { api } from '../lib/api'
import type {
  CreateExpenseInput,
  ExpenseCategory,
  ExpenseSummary,
  ListExpensesResponse,
  PaymentMethod,
} from '../lib/types'

const CATEGORY_LABEL: Record<ExpenseCategory, string> = {
  EXPENSE_CATEGORY_UNSPECIFIED: '미분류',
  EXPENSE_CATEGORY_TRANSPORT: '교통',
  EXPENSE_CATEGORY_FOOD: '식사',
  EXPENSE_CATEGORY_LODGING: '숙박',
  EXPENSE_CATEGORY_ACTIVITY: '체험',
  EXPENSE_CATEGORY_SHOPPING: '쇼핑',
  EXPENSE_CATEGORY_OTHER: '기타',
}

const CATEGORIES: ExpenseCategory[] = [
  'EXPENSE_CATEGORY_TRANSPORT',
  'EXPENSE_CATEGORY_FOOD',
  'EXPENSE_CATEGORY_LODGING',
  'EXPENSE_CATEGORY_ACTIVITY',
  'EXPENSE_CATEGORY_SHOPPING',
  'EXPENSE_CATEGORY_OTHER',
]

const PAYMENTS: PaymentMethod[] = [
  'PAYMENT_METHOD_CASH',
  'PAYMENT_METHOD_CARD',
  'PAYMENT_METHOD_TRANSFER',
]

const PALETTE = ['#f9a8d4', '#fb7185', '#fda4af', '#fbcfe8', '#f472b6', '#ec4899', '#be185d']

export function ExpensesPage() {
  const { tripId = '' } = useParams()
  const qc = useQueryClient()
  const [open, setOpen] = useState(false)

  const list = useQuery({
    queryKey: ['expenses', tripId],
    queryFn: () => api.get<ListExpensesResponse>(`/trips/${tripId}/expenses`),
    enabled: !!tripId,
  })
  const summary = useQuery({
    queryKey: ['expense-summary', tripId],
    queryFn: () => api.get<ExpenseSummary>(`/trips/${tripId}/expense-summary`),
    enabled: !!tripId,
  })

  const createMut = useMutation({
    mutationFn: (input: CreateExpenseInput) => api.post(`/trips/${tripId}/expenses`, input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['expenses', tripId] })
      qc.invalidateQueries({ queryKey: ['expense-summary', tripId] })
      setOpen(false)
    },
  })

  const deleteMut = useMutation({
    mutationFn: (id: string) => api.del(`/expenses/${id}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['expenses', tripId] })
      qc.invalidateQueries({ queryKey: ['expense-summary', tripId] })
    },
  })

  const chartData = (summary.data?.byCategory ?? []).map((c) => ({
    name: CATEGORY_LABEL[c.category] ?? c.category,
    value: Number(c.total?.amount ?? 0),
    currency: c.total?.currency ?? '',
  }))

  return (
    <section className="space-y-4">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행 상세
      </Link>
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-slate-800">지출</h1>
        <button
          onClick={() => setOpen(true)}
          className="flex items-center gap-1 rounded-md bg-sakura-500 px-3 py-1.5 text-white hover:bg-sakura-600"
        >
          <Plus className="h-4 w-4" /> 지출 추가
        </button>
      </div>

      {summary.data?.grandTotal && (
        <div className="rounded-2xl bg-white p-4 shadow-sm">
          <p className="text-sm text-slate-500">총 지출</p>
          <p className="text-2xl font-bold text-sakura-700">
            {summary.data.grandTotal.amount} {summary.data.grandTotal.currency}
          </p>
          {chartData.length > 0 && (
            <div className="h-64">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={chartData}
                    dataKey="value"
                    nameKey="name"
                    innerRadius={50}
                    outerRadius={90}
                    paddingAngle={2}
                  >
                    {chartData.map((_, i) => (
                      <Cell key={i} fill={PALETTE[i % PALETTE.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            </div>
          )}
          {(summary.data.byDay?.length ?? 0) > 0 && (
            <div className="mt-4">
              <h3 className="text-sm font-medium text-slate-700 mb-2">일자별 합계</h3>
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-slate-100 text-left text-slate-500">
                      <th className="py-1 pr-4">날짜</th>
                      <th className="py-1 text-right">금액</th>
                    </tr>
                  </thead>
                  <tbody>
                    {summary.data.byDay.map((d) => (
                      <tr key={d.date} className="border-b border-slate-50">
                        <td className="py-1.5 pr-4 text-slate-700">{d.date}</td>
                        <td className="py-1.5 text-right font-medium text-slate-800">
                          {d.total.amount} {d.total.currency}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}

      <ul className="space-y-2">
        {list.data?.expenses?.map((e) => (
          <li
            key={e.id}
            className="flex items-center justify-between rounded-xl bg-white p-4 shadow-sm"
          >
            <div>
              <p className="text-xs text-sakura-600">{CATEGORY_LABEL[e.category] ?? e.category}</p>
              <p className="font-bold text-slate-800">
                {e.amount.amount} {e.amount.currency}
              </p>
              {e.description && <p className="text-sm text-slate-600">{e.description}</p>}
            </div>
            <button
              onClick={() => deleteMut.mutate(e.id)}
              className="rounded-md p-1 text-slate-400 hover:bg-red-50 hover:text-red-500"
              aria-label="삭제"
            >
              <Trash2 className="h-4 w-4" />
            </button>
          </li>
        ))}
        {list.data?.expenses?.length === 0 && (
          <li className="rounded-xl bg-white p-6 text-center text-slate-500 shadow-sm">
            아직 지출 기록이 없습니다.
          </li>
        )}
      </ul>

      {open && (
        <CreateExpenseModal
          onClose={() => setOpen(false)}
          onSubmit={(v) => createMut.mutate(v)}
          pending={createMut.isPending}
        />
      )}
    </section>
  )
}

function CreateExpenseModal({
  onClose,
  onSubmit,
  pending,
}: {
  onClose: () => void
  onSubmit: (v: CreateExpenseInput) => void
  pending: boolean
}) {
  const [form, setForm] = useState({
    category: 'EXPENSE_CATEGORY_FOOD' as ExpenseCategory,
    amount: '',
    currency: 'JPY',
    paymentMethod: 'PAYMENT_METHOD_CARD' as PaymentMethod,
    description: '',
  })
  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/30 p-4">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
        <h2 className="text-lg font-bold text-slate-800">지출 추가</h2>
        <form
          className="mt-4 space-y-3"
          onSubmit={(e) => {
            e.preventDefault()
            onSubmit({
              category: form.category,
              amount: { currency: form.currency, amount: Number(form.amount) || 0 },
              paymentMethod: form.paymentMethod,
              description: form.description || undefined,
            })
          }}
        >
          <Select
            label="분류"
            value={form.category}
            onChange={(v) => setForm({ ...form, category: v as ExpenseCategory })}
            options={CATEGORIES.map((c) => ({ value: c, label: CATEGORY_LABEL[c] }))}
          />
          <div className="grid grid-cols-2 gap-2">
            <Field
              label="금액"
              type="number"
              required
              value={form.amount}
              onChange={(v) => setForm({ ...form, amount: v })}
            />
            <Field
              label="통화"
              value={form.currency}
              onChange={(v) => setForm({ ...form, currency: v.toUpperCase() })}
            />
          </div>
          <Select
            label="결제 수단"
            value={form.paymentMethod}
            onChange={(v) => setForm({ ...form, paymentMethod: v as PaymentMethod })}
            options={PAYMENTS.map((p) => ({
              value: p,
              label: p.replace('PAYMENT_METHOD_', ''),
            }))}
          />
          <Field
            label="메모"
            value={form.description}
            onChange={(v) => setForm({ ...form, description: v })}
          />
          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={onClose} className="rounded-md border px-3 py-1.5">
              취소
            </button>
            <button
              type="submit"
              disabled={pending || !form.amount}
              className="rounded-md bg-sakura-500 px-3 py-1.5 text-white disabled:opacity-50"
            >
              {pending ? '추가 중…' : '추가'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

function Field({
  label,
  value,
  onChange,
  type = 'text',
  required,
}: {
  label: string
  value: string
  onChange: (v: string) => void
  type?: string
  required?: boolean
}) {
  return (
    <label className="block">
      <span className="text-sm text-slate-700">
        {label}
        {required && <span className="text-sakura-500"> *</span>}
      </span>
      <input
        type={type}
        value={value}
        required={required}
        onChange={(e) => onChange(e.target.value)}
        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 focus:border-sakura-400 focus:outline-none"
      />
    </label>
  )
}

function Select({
  label,
  value,
  onChange,
  options,
}: {
  label: string
  value: string
  onChange: (v: string) => void
  options: { value: string; label: string }[]
}) {
  return (
    <label className="block">
      <span className="text-sm text-slate-700">{label}</span>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 focus:border-sakura-400 focus:outline-none"
      >
        {options.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
    </label>
  )
}
