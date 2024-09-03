import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'error',
    loadComponent: () => import('./error/error.component').then(m => m.ErrorComponent)
  },
  {
    path: 'article-1',
    loadComponent: () => import('./article-1/article-1.component').then(m => m.Article1Component)
  },
  {
    path: 'article-2',
    loadComponent: () => import('./article-2/article-2.component').then(m => m.Article2Component)
  },
]
