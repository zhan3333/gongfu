import { Component } from '@angular/core';
import { CommonModule, NgOptimizedImage } from '@angular/common';
import { ImageComponent } from '../image/image.component';

@Component({
  selector: 'app-article-2',
  standalone: true,
  imports: [CommonModule, ImageComponent, NgOptimizedImage],
  templateUrl: './article-2.component.html',
})
export class Article2Component {

}
