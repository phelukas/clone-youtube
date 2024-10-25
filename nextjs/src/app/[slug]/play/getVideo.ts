import { VideoModel } from "../../../components/models";

export async function getVideo(slug: string): Promise<VideoModel> {
  const response = await fetch(`http://localhost:8000/api/videos/${slug}`, {
    //cache: "no-cache",
    next: {
    //     revalidate: 60 * 5
        tags: [`video-${slug}`]
    }
  });
  return response.json();
}

//revalidate on demand
// /admin/videos -> django -> http -> next.js -> revalidate