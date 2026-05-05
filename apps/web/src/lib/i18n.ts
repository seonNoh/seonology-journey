import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

const resources = {
  ko: {
    translation: {
      common: {
        back: '뒤로',
        save: '저장',
        cancel: '취소',
        delete: '삭제',
        search: '검색',
        searching: '검색 중…',
        loading: '로딩 중…',
        error: '오류가 발생했습니다',
      },
      trip: {
        title: '여행',
        create: '새 여행',
        schedule: '일정 추가',
        day: '일차',
      },
      nearby: {
        title: '주변 검색',
        radius: '반경 (m)',
        placeType: '장소 유형',
        restaurant: '식당',
        cafe: '카페',
        attraction: '관광지',
        convenience: '편의점',
        hotel: '호텔',
      },
      transit: {
        title: '교통 검색',
        origin: '출발지 (lat,lng)',
        destination: '도착지 (lat,lng)',
        departureTime: '출발 시간 (선택)',
        searchRoute: '경로 검색',
        route: '경로',
      },
      share: {
        title: '공유 링크',
        permission: '권한',
        expiry: '만료일',
        create: '링크 생성',
        copy: '복사',
        copied: '복사됨',
        viewer: '뷰어',
        editor: '편집자',
      },
      media: {
        title: '사진/미디어',
        upload: '업로드',
        uploading: '업로드 중…',
      },
      a11y: {
        skipToContent: '본문으로 건너뛰기',
        openMenu: '메뉴 열기',
        closeMenu: '메뉴 닫기',
      },
    },
  },
  ja: {
    translation: {
      common: {
        back: '戻る',
        save: '保存',
        cancel: 'キャンセル',
        delete: '削除',
        search: '検索',
        searching: '検索中…',
        loading: '読み込み中…',
        error: 'エラーが発生しました',
      },
      trip: {
        title: '旅行',
        create: '新規旅行',
        schedule: 'スケジュール追加',
        day: '日目',
      },
      nearby: {
        title: '周辺検索',
        radius: '半径 (m)',
        placeType: '場所の種類',
        restaurant: 'レストラン',
        cafe: 'カフェ',
        attraction: '観光地',
        convenience: 'コンビニ',
        hotel: 'ホテル',
      },
      transit: {
        title: '交通検索',
        origin: '出発地 (lat,lng)',
        destination: '目的地 (lat,lng)',
        departureTime: '出発時間 (任意)',
        searchRoute: '経路検索',
        route: '経路',
      },
      share: {
        title: '共有リンク',
        permission: '権限',
        expiry: '有効期限',
        create: 'リンク作成',
        copy: 'コピー',
        copied: 'コピー済み',
        viewer: '閲覧者',
        editor: '編集者',
      },
      media: {
        title: '写真/メディア',
        upload: 'アップロード',
        uploading: 'アップロード中…',
      },
      a11y: {
        skipToContent: '本文へスキップ',
        openMenu: 'メニューを開く',
        closeMenu: 'メニューを閉じる',
      },
    },
  },
  en: {
    translation: {
      common: {
        back: 'Back',
        save: 'Save',
        cancel: 'Cancel',
        delete: 'Delete',
        search: 'Search',
        searching: 'Searching…',
        loading: 'Loading…',
        error: 'An error occurred',
      },
      trip: {
        title: 'Trip',
        create: 'New Trip',
        schedule: 'Add Schedule',
        day: 'Day',
      },
      nearby: {
        title: 'Nearby Search',
        radius: 'Radius (m)',
        placeType: 'Place Type',
        restaurant: 'Restaurant',
        cafe: 'Cafe',
        attraction: 'Attraction',
        convenience: 'Convenience Store',
        hotel: 'Hotel',
      },
      transit: {
        title: 'Transit Search',
        origin: 'Origin (lat,lng)',
        destination: 'Destination (lat,lng)',
        departureTime: 'Departure Time (optional)',
        searchRoute: 'Search Route',
        route: 'Route',
      },
      share: {
        title: 'Share Link',
        permission: 'Permission',
        expiry: 'Expiry',
        create: 'Create Link',
        copy: 'Copy',
        copied: 'Copied',
        viewer: 'Viewer',
        editor: 'Editor',
      },
      media: {
        title: 'Photos/Media',
        upload: 'Upload',
        uploading: 'Uploading…',
      },
      a11y: {
        skipToContent: 'Skip to content',
        openMenu: 'Open menu',
        closeMenu: 'Close menu',
      },
    },
  },
}

i18n.use(initReactI18next).init({
  resources,
  lng: navigator.language.startsWith('ja')
    ? 'ja'
    : navigator.language.startsWith('en')
      ? 'en'
      : 'ko',
  fallbackLng: 'ko',
  interpolation: { escapeValue: false },
})

export default i18n
