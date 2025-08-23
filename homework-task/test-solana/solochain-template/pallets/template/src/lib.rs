#![cfg_attr(not(feature = "std"), no_std)]

pub use pallet::*;

#[frame_support::pallet]
pub mod pallet {
	use frame_support::pallet_prelude::*;
	use frame_system::pallet_prelude::*;

	#[pallet::pallet]
	pub struct Pallet<T>(_);

	/// 配置接口，用于定义 Pallet 运行所需的类型和参数。
	#[pallet::config]
	pub trait Config: frame_system::Config {
		/// 因为这个 Pallet 会发出事件，所以它需要 `Event` 类型。
		type RuntimeEvent: From<Event<Self>> + IsType<<Self as frame_system::Config>::RuntimeEvent>;
	}

	// [新增] 定义资产存储。
	// 使用 `StorageMap` 来存储从一个唯一的资产 ID (`u32`) 到其所有者账户 (`T::AccountId`) 的映射。
	#[pallet::storage]
	#[pallet::getter(fn asset_owner)]
	pub type AssetOwner<T: Config> = StorageMap<_, Blake2_128Concat, u32, T::AccountId>;

	/// Pallet 的事件。
	#[pallet::event]
	#[pallet::generate_deposit(pub(super) fn deposit_event)]
	pub enum Event<T: Config> {
		// [新增] 资产成功注册的事件。
		// 参数: [资产ID, 注册者]
		AssetRegistered { asset_id: u32, owner: T::AccountId },
		// [新增] 资产成功转移的事件。
		// 参数: [资产ID, 旧所有者, 新所有者]
		AssetTransferred { asset_id: u32, from: T::AccountId, to: T::AccountId },
	}

	/// Pallet 的错误类型。
	#[pallet::error]
	pub enum Error<T> {
		// [新增] 当尝试注册一个已经存在的资产ID时返回此错误。
		AssetIdInUse,
		// [新增] 当一个用户尝试转移不属于他们的资产时返回此错误。
		NotTheOwner,
		// [新增] 当尝试转移一个不存在的资产时返回此错误。
		AssetNotFound,
	}

	/// Pallet 的可调用函数 (Extrinsics)。
	#[pallet::call]
	impl<T: Config> Pallet<T> {
		// [新增] 注册一个新资产。
		// `origin`: 交易的发送方。
		// `asset_id`: 用户希望注册的资产ID。
		#[pallet::call_index(0)]
		#[pallet::weight(10_000)] // 为这个操作指定一个固定的权重（手续费）。
		pub fn register_asset(origin: OriginFor<T>, asset_id: u32) -> DispatchResult {
			// 确保交易是由一个已签名的账户发起的，并获取该账户ID。
			let who = ensure_signed(origin)?;

			// 检查这个资产ID是否已经被注册。
			ensure!(!AssetOwner::<T>::contains_key(&asset_id), Error::<T>::AssetIdInUse);

			// 将新的资产ID和所有者存入链上存储。
			AssetOwner::<T>::insert(&asset_id, &who);

			// 发出一个事件，通知外界资产已成功注册。
			Self::deposit_event(Event::AssetRegistered { asset_id, owner: who });

			// 返回成功。
			Ok(())
		}

		// [新增] 转移一个已存在的资产。
		// `origin`: 交易的发送方。
		// `asset_id`: 要转移的资产ID。
		// `to`: 接收资产的新所有者。
		#[pallet::call_index(1)]
		#[pallet::weight(10_000)]
		pub fn transfer_asset(origin: OriginFor<T>, asset_id: u32, to: T::AccountId) -> DispatchResult {
			// 确保交易是由一个已签名的账户发起的。
			let who = ensure_signed(origin)?;

			// 检查资产是否存在，如果存在则获取其当前所有者。
			let owner = AssetOwner::<T>::get(&asset_id).ok_or(Error::<T>::AssetNotFound)?;

			// 确保交易的发送方就是该资产的当前所有者。
			ensure!(who == owner, Error::<T>::NotTheOwner);

			// 更新资产的所有者为新的接收者。
			AssetOwner::<T>::insert(&asset_id, &to);

			// 发出一个事件，通知外界资产已成功转移。
			Self::deposit_event(Event::AssetTransferred { asset_id, from: who, to });

			// 返回成功。
			Ok(())
		}
	}
}
