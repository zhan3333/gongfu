import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'article-1',
    loadComponent: () => import('./article-1/article-1.component').then(m => m.Article1Component)
  },
]
