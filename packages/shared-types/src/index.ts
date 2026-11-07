export interface HealthResponse {
  status: string;
  service: string;
}

export interface ApiMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface ApiErrorDetail {
  field?: string;
  message: string;
  [key: string]: unknown;
}

export interface ApiError {
  code: string;
  details?: ApiErrorDetail[];
}

export interface ApiResponse<T = unknown> {
  success?: boolean;
  message?: string;
  data?: T;
  meta?: ApiMeta;
  error?: string | ApiError;
}

export type UserRole = 'buyer' | 'seller';

export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  createdAt?: string;
  updatedAt?: string;
}

export interface LoginResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}

export interface ProductSeller {
  id: string;
  name: string;
  email?: string;
}

export interface Product {
  id: string;
  sellerId?: string;
  seller_id?: string;
  name: string;
  description: string;
  price: number;
  stock: number;
  imageUrl?: string | null;
  image_url?: string | null;
  isActive?: boolean;
  is_active?: boolean;
  seller?: ProductSeller;
  createdAt?: string;
  updatedAt?: string;
  created_at?: string;
  updated_at?: string;
}

export interface CartItem {
  id: string;
  cartId?: string;
  cart_id?: string;
  productId: string;
  product_id?: string;
  product: Product;
  quantity: number;
  subtotal: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface Cart {
  id: string;
  buyerId?: string;
  buyer_id?: string;
  items: CartItem[];
  total_items: number;
  total_amount: number;
}

export type OrderStatus =
  | 'pending'
  | 'confirmed'
  | 'processing'
  | 'shipped'
  | 'completed'
  | 'cancelled'
  | 'rejected';

export type PaymentStatus =
  | 'unpaid'
  | 'paid'
  | 'expired'
  | 'failed'
  | 'refunded';

export interface Payment {
  id: string;
  order_id: string;
  invoice_number: string;
  amount: number;
  payment_method?: string;
  channel?: string;
  status: 'PENDING' | 'SUCCESS' | 'EXPIRED' | 'FAILED';
  doku_reference_id?: string;
  created_at: string;
  updated_at: string;
}

export interface PayOrderResponse {
  order_id: string;
  invoice_number: string;
  amount: number;
  payment_url: string;
  payment_status: PaymentStatus;
}

export interface ShippingAddress {
  recipient_name: string;
  phone: string;
  address: string;
  city: string;
  province?: string;
  postal_code: string;
}

export interface OrderItem {
  id?: string;
  productId: string;
  product_id?: string;
  productName: string;
  product_name?: string;
  unitPrice: number;
  unit_price?: number;
  quantity: number;
  subtotal: number;
}

export interface Order {
  id: string;
  buyerId?: string;
  buyer_id?: string;
  sellerId?: string;
  seller_id?: string;
  totalAmount: number;
  total_amount?: number;
  status: OrderStatus;
  payment_status?: PaymentStatus;
  paymentStatus?: PaymentStatus;
  invoice_number?: string | null;
  invoiceNumber?: string | null;
  payment_url?: string | null;
  paymentUrl?: string | null;
  shippingAddress: ShippingAddress;
  shipping_address?: ShippingAddress;
  notes?: string | null;
  items: OrderItem[];
  buyer?: Pick<User, 'id' | 'name' | 'email'>;
  seller?: Pick<User, 'id' | 'name' | 'email'>;
  createdAt?: string;
  updatedAt?: string;
  created_at?: string;
  updated_at?: string;
}
